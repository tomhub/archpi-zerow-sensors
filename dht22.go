package main

import (
	"fmt"
	dht "github.com/d2r2/go-dht"
    logger "github.com/d2r2/go-logger"
    
    "time"
    "github.com/influxdata/influxdb-client-go"
    "math"
)

var lg = logger.NewPackageLogger("main",
	logger.DebugLevel,
// 	logger.InfoLevel,
)

// cross-build:
// GOOS=linux GOARCH=arm go build dht22.go

func main() {
    defer logger.FinalizeLogger()
    
    userName := "INFLUXDB_USERNAME"
    //it seems password with spaces is not handled properly (yet)
    password := "INFLUXDB_PASSWORD"
    
//     logger.ChangePackageLogLevel("dht", logger.InfoLevel)
    logger.ChangePackageLogLevel("dht", logger.ErrorLevel) // only errors and fatals
    
    
    // Create a new client using an InfluxDB server base URL and an authentication token
    // For authentication token supply a string in the form: "username:password" as a token. Set empty value for an unauthenticated server
    client := influxdb2.NewClientWithOptions("http://127.0.0.1:8086", fmt.Sprintf("%s:%s",userName, password),
                                             influxdb2.DefaultOptions().SetBatchSize(200).SetFlushInterval(60000).SetRetryBufferLimit(100000).SetUseGZip(true))
    // Get the blocking write client
    // Supply a string in the form database/retention-policy as a bucket. Skip retention policy for the default one, use just a database name (without the slash character)
    // Org name is not used
    writeAPI := client.WriteAPI("", "INFLUXDB_DATABASE")
    // flush when destroyed
    defer writeAPI.Flush();
    // Close client when destroyed
    defer client.Close()

    for {
    
        // Read DHT11 sensor data from pin 4, retrying 10 times in case of failure.
        // You may enable "boost GPIO performance" parameter, if your device is old
        // as Raspberry PI 1 (this will require root privileges). You can switch off
        // "boost GPIO performance" parameter for old devices, but it may increase
        // retry attempts. Play with this parameter.
        // Note: "boost GPIO performance" parameter is not work anymore from some
        // specific Go release. Never put true value here.
        t, h, retried, err :=
            dht.ReadDHTxxWithRetry(dht.DHT22, 17, false, 10)
        if err != nil {
            lg.Fatal(err)
        } else {
            p := influxdb2.NewPoint("INFLUXDB_MEASUREMENT",
                map[string]string{
                    "TAG1": "This tag"},
                map[string]interface{}{"temperature": math.Round(float64(t)*100)/100,
                    "humidity": math.Round(float64(h)*100)/100,
                },
                time.Now().UTC())
            writeAPI.WritePoint(p)
            time.Sleep(5*time.Second)
        }
        _ = retried
        // Print temperature and humidity
//         fmt.Printf("Temperature = %v*C, Humidity = %v%% (retried %d times)\n",
//             temperature, humidity, retried)
    }
}
