package main

import (
    "fmt"
    "github.com/yryz/ds18b20"
    "time"
    "github.com/influxdata/influxdb-client-go"
    "math"
)


func main() {
    userName := "INFLUXDB_USERNAME"
    password := "INFLUXDB_PASSWORD"
    
    sensors, err := ds18b20.Sensors()
    if err != nil {
        panic(err)
    }

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
    var t float64
    for {
        for _, sensor := range sensors {
            t, err = ds18b20.Temperature(sensor)
            if err == nil {
                p := influxdb2.NewPoint("pi-ds18b20",
                    map[string]string{
                        "location": "fishtank",
                        "sensor_id": sensor[3:len(sensor)]},
                    map[string]interface{}{"temperature": math.Round(t*100)/100,},
                    time.Now().UTC())
//                 fmt.Printf("sensor: %s temperature: %.2f°C\n", sensor, t)
                writeAPI.WritePoint(p)
            } else {
                fmt.Printf("DS18B20: Read error: %s\n", err.Error())
            }
        }
        time.Sleep(5*time.Second)
    }
}



