# Test of DataDog

Querying various metrics in MTR and pushing metrics to the DataDog service
using the DataDog statsd package.  The production approach would be different, 
probably pushing the metrics to DataDog as they're being handled in mtr-api.