package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	dynatraceConfigV1 "github.com/dynatrace-ace/dynatrace-go-api-client/api/v1/config/dynatrace"
	// dynatraceEnvironmentV1 "github.com/dynatrace-ace/dynatrace-go-api-client/api/v1/environment/dynatrace"
	// dynatraceEnvironmentV2 "github.com/dynatrace-ace/dynatrace-go-api-client/api/v2/environment/dynatrace"
)

var OUTPUT_PATHS map[string]string

// Process command line arguments
var dtEnvURL = os.Args[1]
var apiToken = os.Args[2]
var subdir = os.Args[3]

func main() {
	// fmt.Println("func main")

	fmt.Println("URL:                            ", dtEnvURL)
	fmt.Println("Token (1st 15 characters only): ", apiToken[:15])
	fmt.Println("Subdirectory:                   ", subdir)

	// Load Output Paths for writing entity types
	OUTPUT_PATHS = make(map[string]string)
	OUTPUT_PATHS["Service Request Naming"] = subdir + "/api/config/v1/service/requestNaming"
	OUTPUT_PATHS["Dashboard"] = subdir + "/api/config/v1/dashboards"

	parsedDTUrl, err := url.Parse(dtEnvURL)
	if err != nil {
		fmt.Println(err)
	}
	// authEnvironmentV1 := context.WithValue(
	// 	context.Background(),
	// 	dynatraceEnvironmentV1.ContextAPIKeys,
	// 	map[string]dynatraceEnvironmentV1.APIKey{
	// 		"Api-Token": {
	// 			Key:    apiToken,
	// 			Prefix: "Api-Token",
	// 		},
	// 	},
	// )

	// authEnvironmentV1 = context.WithValue(authEnvironmentV1, dynatraceEnvironmentV1.ContextServerVariables, map[string]string{
	// 	"name":     string(parsedDTUrl.Host + parsedDTUrl.Path),
	// 	"protocol": string(parsedDTUrl.Scheme),
	// })

	// EnvironmentV1 := dynatraceEnvironmentV1.NewConfiguration()
	// dynatraceEnvironmentClientV1 := dynatraceEnvironmentV1.NewAPIClient(EnvironmentV1)

	// authEnvironmentV2 := context.WithValue(
	// 	context.Background(),
	// 	dynatraceEnvironmentV2.ContextAPIKeys,
	// 	map[string]dynatraceEnvironmentV2.APIKey{
	// 		"Api-Token": {
	// 			Key:    apiToken,
	// 			Prefix: "Api-Token",
	// 		},
	// 	},
	// )

	// authEnvironmentV2 = context.WithValue(authEnvironmentV2, dynatraceEnvironmentV2.ContextServerVariables, map[string]string{
	// 	"name":     string(parsedDTUrl.Host + parsedDTUrl.Path),
	// 	"protocol": string(parsedDTUrl.Scheme),
	// })

	// EnvironmentV2 := dynatraceEnvironmentV2.NewConfiguration()
	// dynatraceEnvironmentClientV2 := dynatraceEnvironmentV2.NewAPIClient(EnvironmentV2)

	authConfigV1 := context.WithValue(
		context.Background(),
		dynatraceConfigV1.ContextAPIKeys,
		map[string]dynatraceConfigV1.APIKey{
			"Api-Token": {
				Key:    apiToken,
				Prefix: "Api-Token",
			},
		},
	)

	authConfigV1 = context.WithValue(authConfigV1, dynatraceConfigV1.ContextServerVariables, map[string]string{
		"name":     string(parsedDTUrl.Host + parsedDTUrl.Path),
		"protocol": string(parsedDTUrl.Scheme),
	})

	configV1 := dynatraceConfigV1.NewConfiguration()
	dynatraceConfigClientV1 := dynatraceConfigV1.NewAPIClient(configV1)

	getServiceRequestNamingList(authConfigV1, dynatraceConfigClientV1)
	getDashboardList(authConfigV1, dynatraceConfigClientV1)
	// getHostList(authEnvironmentV1, dynatraceEnvironmentClientV1)
	// getServiceList(authEnvironmentV1, dynatraceEnvironmentClientV1)
	// getServiceMethodList(authEnvironmentV2, dynatraceEnvironmentClientV2)
}

// func getHostList(authEnvironmentV1 context.Context, dynatraceEnvironmentClientV1 *dynatraceEnvironmentV1.APIClient) {
// 	fmt.Println("func getHostList")
// 	hostList, _, err := dynatraceEnvironmentClientV1.TopologySmartscapeHostApi.GetHosts(authEnvironmentV1).IncludeDetails(false).Execute()
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	hostResponse, err := json.MarshalIndent(hostList, "", "    ")
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	log.Printf("[DEBUG] Host List: \n %s \n \n", hostResponse)
// }

// func getServiceList(authEnvironmentV1 context.Context, dynatraceEnvironmentClientV1 *dynatraceEnvironmentV1.APIClient) {
// 	fmt.Println("func getServiceList")
// 	serviceList, _, err := dynatraceEnvironmentClientV1.TopologySmartscapeServiceApi.GetServices(authEnvironmentV1).Execute()
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	serviceResponse, err := json.MarshalIndent(serviceList, "", "    ")
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	log.Printf("[DEBUG] Service List: \n %s \n \n", serviceResponse)

// 	// for _, service := range serviceList {
// 	// 	fmt.Println(service.GetEntityId(), service.GetDisplayName())
// 	// }
// }

// func getServiceMethodList(authEnvironmentV2 context.Context, dynatraceEnvironmentClientV2 *dynatraceEnvironmentV2.APIClient) {
// 	fmt.Println("func getServiceMethodList")
// 	serviceMethod, _, err := dynatraceEnvironmentClientV2.MonitoredEntitiesApi.GetEntities(authEnvironmentV2).EntitySelector("type(\"SERVICE_METHOD\")").Execute()
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	serviceMethodResponse, err := json.MarshalIndent(serviceMethod, "", "    ")
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	log.Printf("[DEBUG] Service Method List: \n %s \n \n", serviceMethodResponse)

// 	// entities := *serviceMethod.Entities

// 	// for _, entity := range entities {
// 	// 	fmt.Println(entity.GetEntityId(), entity.GetDisplayName())
// 	// }
// }

func getDashboardList(authConfigV1 context.Context, dynatraceConfigClientV1 *dynatraceConfigV1.APIClient) {
	//fmt.Println("func getDashboardList")
	dashboardList, _, err := dynatraceConfigClientV1.DashboardsApi.GetDashboardStubsList(authConfigV1).Execute()
	if err != nil {
		fmt.Println(err)
	}

	// dashboardResponse, err := json.MarshalIndent(dashboardList, "", "    ")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	//log.Printf("[DEBUG] Dashboard List: \n %s \n \n", dashboardResponse)

	for _, dashboard := range dashboardList.GetDashboards() {
		id := dashboard.GetId()
		//getDashboard(authConfigV1, dynatraceConfigClientV1, id)
		dashboard := getDashboardDirect(dtEnvURL, apiToken, id)
		writeEntity("Dashboard", id, dashboard)
	}
}

func getServiceRequestNamingList(authConfigV1 context.Context, dynatraceConfigClientV1 *dynatraceConfigV1.APIClient) {
	//fmt.Println("func getServiceRequestNamingList")
	serviceRequestNamingList, _, err := dynatraceConfigClientV1.ServiceRequestNamingApi.ListRequestNaming(authConfigV1).Execute()
	if err != nil {
		fmt.Println(err)
	}

	// serviceRequestNamingResponse, err := json.MarshalIndent(serviceRequestNamingList, "", "    ")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	//log.Printf("[DEBUG] Service Request Naming List: \n %s \n \n", serviceRequestNamingResponse)

	for _, serviceRequestNaming := range serviceRequestNamingList.GetValues() {
		id := serviceRequestNaming.GetId()
		//name := serviceRequestNaming.GetName()
		//fmt.Println(id, name)
		//getServiceRequestNaming(authConfigV1, dynatraceConfigClientV1, id)
		getServiceRequestNamingDirect(dtEnvURL, apiToken, id)
	}
}

// func getServiceRequestNaming(authConfigV1 context.Context, dynatraceConfigClientV1 *dynatraceConfigV1.APIClient, id string) {
// 	fmt.Println("func getServiceRequestNaming")
// 	serviceRequestNaming, _, err := dynatraceConfigClientV1.ServiceRequestNamingApi.GetRequestNaming(authConfigV1, id).Execute()
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	serviceRequestNamingResponse, err := json.MarshalIndent(serviceRequestNaming, "", "    ")
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	//log.Printf("[DEBUG] Service Request Naming: \n %s \n \n", serviceRequestNamingResponse)
// 	writeEntity("Service Request Naming", id, string(serviceRequestNamingResponse))
// }

func getServiceRequestNamingDirect(dtEnvURL string, token string, id string) {
	//fmt.Println("func getServiceRequestNamingDirect")
	localVarPath := dtEnvURL + "/api/config/v1/service/requestNaming/" + id

	client := &http.Client{}

	req, err := http.NewRequest("GET", localVarPath, nil)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Authorization", "Api-Token "+token)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(req)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		//panic(err)
	}

	defer resp.Body.Close()
	body, err2 := io.ReadAll(resp.Body)

	if err2 != nil {
		fmt.Println(err2)
		//panic(err2)
	}

	//fmt.Println(string(body))

	var out bytes.Buffer
	err3 := json.Indent(&out, body, "", "    ")

	if err3 != nil {
		fmt.Println(err3)
		//panic(err3)
	}

	writeEntity("Service Request Naming", id, out.String())
}

func getDashboard(authConfigV1 context.Context, dynatraceConfigClientV1 *dynatraceConfigV1.APIClient, id string) {
	//fmt.Println("func getDashboard")
	dashboard, _, err := dynatraceConfigClientV1.DashboardsApi.GetDashboard(authConfigV1, id).Execute()
	if err != nil {
		fmt.Println(err)
	}

	dashboardResponse, err := json.MarshalIndent(dashboard, "", "    ")
	if err != nil {
		fmt.Println(err)
	}

	//log.Printf("[DEBUG] Dashboard: \n %s \n \n", dashboardResponse)
	writeEntity("Dashboard", id, string(dashboardResponse))
}

func getDashboardDirect(dtEnvURL string, token string, id string) (dashboard string) {
	localVarPath := dtEnvURL + "/api/config/v1/dashboards/" + id

	client := &http.Client{}

	req, err := http.NewRequest("GET", localVarPath, nil)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Authorization", "Api-Token "+token)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(req)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	defer resp.Body.Close()
	body, err2 := io.ReadAll(resp.Body)

	if err2 != nil {
		fmt.Println(err2)
		panic(err2)
	}

	var out bytes.Buffer
	err3 := json.Indent(&out, body, "", "    ")

	if err3 != nil {
		fmt.Println(err3)
		panic(err3)
	}

	// if id == "bbbbbbbb-a001-a008-0000-000000000005" {
	// 	fmt.Println(out.String())
	// }

	return out.String()
}

func writeEntity(entity_type string, id string, entity string) {
	//fmt.Println("func writeEntity")
	// log.Printf("Entity Type: %s", entity_type)
	// log.Printf("Entity ID:   %s", id)
	// log.Printf("Entity:   \n %s", entity)

	subdir := OUTPUT_PATHS[entity_type]
	makedir(subdir)
	fname := subdir + "/" + id
	//fmt.Println(fname)

	f, err := os.Create(fname)
	check(err)
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err2 := w.WriteString(entity)
	check(err2)
	w.Flush()
}

func makedir(subdir string) {
	//fmt.Println("func makedir")
	_, err := os.Stat(subdir)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(subdir, 0755)
		if errDir != nil {
			fmt.Println(err)
		}
	}
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}
