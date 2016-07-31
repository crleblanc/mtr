package main

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/GeoNet/mtr/mtrpb"
	"github.com/GeoNet/mtr/ts"
	"github.com/GeoNet/weft"
	"github.com/golang/protobuf/proto"
	"github.com/lib/pq"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var statusTooManyRequests = weft.Result{Ok: false, Code: http.StatusTooManyRequests, Msg: "Already data for the minute"}

// fieldMetric - table field.metric
type fieldMetric struct {
}

func fieldMetricPut(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	v := r.URL.Query()

	var err error
	var val int
	var t time.Time

	if val, err = strconv.Atoi(v.Get("value")); err != nil {
		return weft.BadRequest("invalid value")
	}

	if t, err = time.Parse(time.RFC3339, v.Get("time")); err != nil {
		return weft.BadRequest("invalid time")
	}

	deviceID := v.Get("deviceID")
	typeID := v.Get("typeID")

	var result sql.Result

	if result, err = db.Exec(`INSERT INTO field.metric(devicePK, typePK, rate_limit, time, value)
				SELECT devicePK, typePK, $3, $4, $5
				FROM field.device, field.type
				WHERE deviceID = $1
				AND typeID = $2`,
		deviceID, typeID, t.Truncate(time.Minute).Unix(), t, int32(val)); err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == errorUniqueViolation {
			return &statusTooManyRequests
		} else {
			return weft.InternalServerError(err)
		}
	}

	var i int64
	if i, err = result.RowsAffected(); err != nil {
		return weft.InternalServerError(err)
	}
	if i != 1 {
		return weft.BadRequest("Didn't create row, check your query parameters exist")
	}

	// Update the summary value if the incoming value is newer.
	if result, err = db.Exec(`UPDATE field.metric_summary SET time = $3, value = $4
				WHERE time < $3
				AND devicePK = (SELECT devicePK FROM field.device WHERE deviceID = $1)
				AND typePK = (SELECT typePK FROM field.type WHERE typeID = $2)`,
		deviceID, typeID, t, int32(val)); err != nil {
		return weft.InternalServerError(err)
	}

	// If no rows change either the value is old or it's the first time we've seen this metric.
	if i, err = result.RowsAffected(); err != nil {
		return weft.InternalServerError(err)
	}
	if i != 1 {
		if result, err = db.Exec(`INSERT INTO field.metric_summary(devicePK, typePK, time, value)
				SELECT devicePK, typePK, $3, $4
				FROM field.device, field.type
				WHERE deviceID = $1
				AND typeID = $2`,
			deviceID, typeID, t, int32(val)); err != nil {
			if err, ok := err.(*pq.Error); ok && err.Code == errorUniqueViolation {
				// incoming value was old
			} else {
				return weft.InternalServerError(err)
			}
		}
	}

	return &weft.StatusOK
}

func fieldMetricDelete(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	v := r.URL.Query()

	deviceID := v.Get("deviceID")
	typeID := v.Get("typeID")

	var err error
	var txn *sql.Tx

	if txn, err = db.Begin(); err != nil {
		return weft.InternalServerError(err)
	}

	for _, table := range []string{"field.metric", "field.metric_summary", "field.metric_tag", "field.threshold"} {
		if _, err = txn.Exec(`DELETE FROM `+table+` WHERE
				devicePK = (SELECT devicePK FROM field.device WHERE deviceID = $1)
				 AND typePK = (SELECT typePK from field.type WHERE typeID = $2)`,
			deviceID, typeID); err != nil {
			txn.Rollback()
			return weft.InternalServerError(err)
		}
	}

	if err = txn.Commit(); err != nil {
		return weft.InternalServerError(err)
	}

	return &weft.StatusOK
}

func fieldMetricSvg(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	f := fieldMetric{}

	v := r.URL.Query()

	switch r.URL.Query().Get("plot") {
	case "", "line":
		resolution := r.URL.Query().Get("resolution")
		if resolution == "" {
			resolution = "minute"
		}
		if res := f.plot(v.Get("deviceID"), v.Get("typeID"), resolution, ts.Line, b); !res.Ok {
			return res
		}
	case "scatter":
		resolution := r.URL.Query().Get("resolution")
		if resolution == "" {
			resolution = "minute"
		}
		if res := f.plot(v.Get("deviceID"), v.Get("typeID"), resolution, ts.Scatter, b); !res.Ok {
			return res
		}
	default:
		if res := f.spark(v.Get("deviceID"), v.Get("typeID"), b); !res.Ok {
			return res
		}
	}

	return &weft.StatusOK
}

func fieldMetricCsv(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	v := r.URL.Query()
	resolution := v.Get("resolution")
	if resolution == "" {
		resolution = "minute"
	}

	deviceID := v.Get("deviceID")
	typeID := v.Get("typeID")
	var err error

	var timeRange []time.Time
	if timeRange, err = parseTimeRange(v); err != nil {
		return weft.InternalServerError(err)
	}

	var devicePK int
	if err = dbR.QueryRow(`SELECT devicePK FROM field.device WHERE deviceID = $1`,
		deviceID).Scan(&devicePK); err != nil {
		if err == sql.ErrNoRows {
			return &weft.NotFound
		}
		return weft.InternalServerError(err)
	}

	var typePK int
	var scale float64
	if err = dbR.QueryRow(`SELECT typePK, scale FROM field.type WHERE typeID = $1`,
		typeID).Scan(&typePK, &scale); err != nil {
		if err == sql.ErrNoRows {
			return &weft.NotFound
		}
		return weft.InternalServerError(err)
	}

	var rows *sql.Rows
	rows, err = queryMetricRows(devicePK, typePK, resolution, timeRange)
	if err != nil {
		return weft.InternalServerError(err)
	}
	defer rows.Close()

	w := csv.NewWriter(b)
	i := 0
	for rows.Next() {

		if i == 0 {
			w.Write([]string{"time", "value"})
		}

		var val float64
		var t time.Time
		if err = rows.Scan(&t, &val); err != nil {
			return weft.InternalServerError(err)
		}

		w.Write([]string{t.Format(DYGRAPH_TIME_FORMAT), fmt.Sprintf("%.2f", val*scale)})
		i++
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return weft.InternalServerError(err)
	}

	return &weft.StatusOK
}

// proto's query is the same as svg. The difference between them is only output mimetype.
func fieldMetricProto(r *http.Request, h http.Header, b *bytes.Buffer) *weft.Result {
	v := r.URL.Query()
	resolution := v.Get("resolution")
	if resolution == "" {
		resolution = "minute"
	}

	deviceID := v.Get("deviceID")
	typeID := v.Get("typeID")
	var err error

	var fmr mtrpb.FieldMetricResult
	fmr.DeviceID = deviceID
	fmr.TypeID = typeID

	var devicePK int
	if err = dbR.QueryRow(`SELECT devicePK FROM field.device WHERE deviceID = $1`,
		deviceID).Scan(&devicePK); err != nil {
		if err == sql.ErrNoRows {
			return &weft.NotFound
		}
		return weft.InternalServerError(err)
	}

	var typePK int

	if err = dbR.QueryRow(`SELECT typePK FROM field.type WHERE typeID = $1`,
		typeID).Scan(&typePK); err != nil {
		if err == sql.ErrNoRows {
			return &weft.NotFound
		}
		return weft.InternalServerError(err)
	}

	if err := dbR.QueryRow(`SELECT lower,upper FROM field.threshold
		WHERE devicePK = $1 AND typePK = $2`,
		devicePK, typePK).Scan(&fmr.Lower, &fmr.Upper); err != nil && err != sql.ErrNoRows {
		return weft.InternalServerError(err)
	}

	var rows *sql.Rows
	rows, err = queryMetricRows(devicePK, typePK, resolution, nil)
	if err != nil {
		return weft.InternalServerError(err)
	}

	defer rows.Close()
	for rows.Next() {
		var fm mtrpb.FieldMetric
		var t time.Time
		if err = rows.Scan(&t, &fm.Value); err != nil {
			return weft.InternalServerError(err)
		}

		fm.Seconds = t.Unix()
		fmr.Result = append(fmr.Result, &fm)
	}

	var by []byte

	if by, err = proto.Marshal(&fmr); err != nil {
		return weft.InternalServerError(err)
	}

	b.Write(by)

	return &weft.StatusOK
}

/*
plot draws an svg plot to b.
Valid values for resolution are 'minute', 'five_minutes', 'hour'.
*/
func (f fieldMetric) plot(deviceID, typeID, resolution string, plotter ts.SVGPlot, b *bytes.Buffer) *weft.Result {
	// we need the devicePK often so read it once.
	var devicePK int
	if err := dbR.QueryRow(`SELECT devicePK FROM field.device WHERE deviceID = $1`,
		deviceID).Scan(&devicePK); err != nil {
		if err == sql.ErrNoRows {
			return &weft.NotFound
		}
		return weft.InternalServerError(err)
	}

	var typePK int
	var scale float64
	var display string

	if err := dbR.QueryRow(`SELECT typePK, scale, display FROM field.type WHERE typeID = $1`,
		typeID).Scan(&typePK, &scale, &display); err != nil {
		if err == sql.ErrNoRows {
			return &weft.NotFound
		}
		return weft.InternalServerError(err)
	}

	var p ts.Plot

	p.SetUnit(display)

	var rows *sql.Rows
	var err error
	var lower, upper int

	if err := dbR.QueryRow(`SELECT lower,upper FROM field.threshold
		WHERE devicePK = $1 AND typePK = $2`,
		devicePK, typePK).Scan(&lower, &upper); err != nil && err != sql.ErrNoRows {
		return weft.InternalServerError(err)
	}

	if !(lower == 0 && upper == 0) {
		p.SetThreshold(float64(lower)*scale, float64(upper)*scale)
	}

	var tags []string

	if rows, err = dbR.Query(`SELECT tag FROM field.metric_tag JOIN mtr.tag USING (tagpk) WHERE
		devicePK = $1 AND typePK = $2
		ORDER BY tag asc`,
		devicePK, typePK); err != nil {
		return weft.InternalServerError(err)
	}

	defer rows.Close()

	for rows.Next() {
		var s string
		if err = rows.Scan(&s); err != nil {
			return weft.InternalServerError(err)
		}
		tags = append(tags, s)
	}

	rows.Close()

	p.SetSubTitle("Tags: " + strings.Join(tags, ","))

	var mod string
	// TODO move into first select for devicePK
	if err = dbR.QueryRow(`SELECT modelid FROM field.device JOIN field.model using (modelpk)
		WHERE devicePK = $1`,
		devicePK).Scan(&mod); err != nil && err != sql.ErrNoRows {
		return weft.InternalServerError(err)
	}

	p.SetTitle(fmt.Sprintf("Device: %s, Model: %s, Metric: %s", deviceID, mod, strings.Title(typeID)))

	switch resolution {
	case "minute":
		p.SetXAxis(time.Now().UTC().Add(time.Hour*-12), time.Now().UTC())
		p.SetXLabel("12 hours")
	case "five_minutes":
		p.SetXAxis(time.Now().UTC().Add(time.Hour*-24*2), time.Now().UTC())
		p.SetXLabel("48 hours")
	case "hour":
		p.SetXAxis(time.Now().UTC().Add(time.Hour*-24*28), time.Now().UTC())
		p.SetXLabel("4 weeks")
	default:
		return weft.BadRequest("invalid resolution")
	}

	rows, err = queryMetricRows(devicePK, typePK, resolution, nil)
	if err != nil {
		return weft.InternalServerError(err)
	}

	defer rows.Close()

	var pts []ts.Point

	for rows.Next() {
		var pt ts.Point
		if err = rows.Scan(&pt.DateTime, &pt.Value); err != nil {
			return weft.InternalServerError(err)
		}
		pt.Value = pt.Value * scale
		pts = append(pts, pt)
	}
	rows.Close()

	// Add the latest value to the plot - this may be different to the average at minute or hour resolution.
	var pt ts.Point

	if err = dbR.QueryRow(`SELECT time, value FROM field.metric WHERE
			devicePK = $1 AND typePK = $2
			ORDER BY time DESC
			LIMIT 1`,
		devicePK, typePK).Scan(&pt.DateTime, &pt.Value); err != nil {
		return weft.InternalServerError(err)
	}

	pt.Value = pt.Value * scale

	pts = append(pts, pt)
	p.SetLatest(pt, "deepskyblue")

	p.AddSeries(ts.Series{Colour: "deepskyblue", Points: pts})

	if err = plotter.Draw(p, b); err != nil {
		return weft.InternalServerError(err)
	}

	return &weft.StatusOK
}

// spark draws an svg spark line to b.
func (f fieldMetric) spark(deviceID, typeID string, b *bytes.Buffer) *weft.Result {
	var p ts.Plot

	p.SetXAxis(time.Now().UTC().Add(time.Hour*-12), time.Now().UTC())

	var err error
	var rows *sql.Rows

	if rows, err = dbR.Query(`SELECT date_trunc('hour', time) + extract(minute from time)::int / 5 * interval '5 min' as t,
		 avg(value) FROM field.metric
		 WHERE devicePK = (SELECT devicePK from field.device WHERE deviceID = $1)
		AND typePK = (SELECT typePK from field.type WHERE typeID = $2)
		AND time > now() - interval '12 hours'
		GROUP BY date_trunc('hour', time) + extract(minute from time)::int / 5 * interval '5 min'
		ORDER BY t ASC`,
		deviceID, typeID); err != nil {
		return weft.InternalServerError(err)
	}
	defer rows.Close()

	var pts []ts.Point

	for rows.Next() {
		var pt ts.Point
		if err = rows.Scan(&pt.DateTime, &pt.Value); err != nil {
			return weft.InternalServerError(err)
		}
		// No need to scale spark data for display.
		pts = append(pts, pt)
	}
	rows.Close()

	p.AddSeries(ts.Series{Colour: "deepskyblue", Points: pts})

	if err = ts.SparkLine.Draw(p, b); err != nil {
		return weft.InternalServerError(err)
	}

	return &weft.StatusOK
}

func queryMetricRows(devicePK, typePK int, resolution string, timeRange []time.Time) (*sql.Rows, error) {
	var err error
	var rows *sql.Rows

	if timeRange == nil {
		if timeRange, err = defaultTimeRange(resolution); err != nil {
			weft.InternalServerError(err)
		}
	}

	switch resolution {
	case "minute":
		rows, err = dbR.Query(`SELECT date_trunc('`+resolution+`',time) as t, avg(value) FROM field.metric WHERE
		devicePK = $1 AND typePK = $2
		AND time >= $3 AND time <= $4
		GROUP BY date_trunc('`+resolution+`',time)
		ORDER BY t ASC`,
			devicePK, typePK, timeRange[0], timeRange[1])
	case "five_minutes":
		rows, err = dbR.Query(`SELECT date_trunc('hour', time) + extract(minute from time)::int / 5 * interval '5 min' as t,
		 avg(value) FROM field.metric WHERE
		devicePK = $1 AND typePK = $2
		AND time >= $3 AND time <= $4
		GROUP BY date_trunc('hour', time) + extract(minute from time)::int / 5 * interval '5 min'
		ORDER BY t ASC`,
			devicePK, typePK, timeRange[0], timeRange[1])
	case "hour":
		rows, err = dbR.Query(`SELECT date_trunc('`+resolution+`',time) as t, avg(value) FROM field.metric WHERE
		devicePK = $1 AND typePK = $2
		AND time >= $3 AND time <= $4
		GROUP BY date_trunc('`+resolution+`',time)
		ORDER BY t ASC`,
			devicePK, typePK, timeRange[0], timeRange[1])
	case "full":
		rows, err = dbR.Query(`SELECT time, value
		FROM field.metric
		WHERE devicePK = $1
		AND typePK = $2
		AND time >= $3 AND time <= $4
		ORDER BY time ASC`,
			devicePK, typePK, timeRange[0], timeRange[1])
	default:
		return nil, errors.New("invalid resolution")
	}
	if err != nil {
		return nil, err
	}
	return rows, nil
}
