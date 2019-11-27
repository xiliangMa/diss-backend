package utils

func ExtractHostInfo(hostInfo []interface{}) map[string]interface{} {
	//var hostMetricPure []map[string]string
	var metricPure = make(map[string]interface{})

	for _, x := range hostInfo {
		hostMetricL1 := x.(map[string]interface{})["system"]
		hostMetricL2 := hostMetricL1.(map[string]interface{})
		for metrictype, mctype := range hostMetricL2 {
			hostMetricL3 := mctype.(map[string]interface{})
			for metricunit, metricval := range hostMetricL3 {
				metricPure[metrictype+"."+metricunit] = metricval.(float64)
			}
		}
	}

	return metricPure
}
