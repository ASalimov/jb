package cmd

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"github.com/gobuffalo/envy"
)

type JobDetails struct {
	Jobs        []struct {
		Name  string `json:"name"`
		URL   string `json:"url"`
	} `json:"jobs"`
	Views       []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"views"`
}

type JHelper struct {
	jobDetails JobDetails
}

type BuildInfo struct {
	Id string `json:"id"`
	Duration int `json:"duration"`
	Building bool `json:"building"`
	Result string `json:"result"`
}
type ParameterDefinitions struct {
	DefaultParameterValue struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"defaultParameterValue"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Type        string `json:"type"`
}

type JobInfo struct {
	NextBuildNumber int `json:"nextBuildNumber"`
	Property              []struct {
		ParameterDefinitions []ParameterDefinitions  `json:"parameterDefinitions,omitempty"`
	} `json:"property"`
}

func (ji *JobInfo)getParameterDefinitions()[]ParameterDefinitions{
	for _,j:= range ji.Property{
		if len(j.ParameterDefinitions)>0{
			return j.ParameterDefinitions
		}
	}
	return []ParameterDefinitions{}
}

func CreateJHelper() JHelper {
	jh := JHelper{}
	jh.getJobsInfo()
	return jh
}
func (jh *JHelper) getJobsInfo() JobDetails {
	code, rsp,_, err:= req("api/json", []byte{})
	if err!=nil{
		panic(err)
	}
	if code!=200{
		panic(errors.New("failed to get job details"))
	}
	var jobDetails JobDetails
	err = json.Unmarshal(rsp, &jobDetails)
	if err!=nil{
		panic(err)
	}
	jh.jobDetails = jobDetails
	return jobDetails

}

func (jh *JHelper) getJobInfo(job string) *JobInfo {
	code, rsp,_, err:= req("job/"+job+"/api/json", []byte{})
	if err!=nil{
		panic(err)
	}
	if code!=200{
		panic("failed to get job details,code"+strconv.Itoa(code)+ ", "+string(rsp))
	}
	var ji JobInfo
	err = json.Unmarshal(rsp, &ji)
	if err!=nil{
		panic (err)
	}
	return &ji
}

func (jh *JHelper) getBuildInfo(job string, id int) (*BuildInfo, error){
	code, rsp,_, err:= req("job/"+job+"/"+strconv.Itoa(id)+"/api/json", []byte{})
	if err!=nil{
		panic(err)
	}
	if code!=200{
		return nil, errors.New("failed to get job details,code"+strconv.Itoa(code)+ ", "+string(rsp))
	}
	var bi BuildInfo
	err = json.Unmarshal(rsp, &bi)
	if err!=nil{
		return nil, err
	}
	return &bi, nil
}

func (jh *JHelper) getLastSuccessfulBuildInfo(job string) (*BuildInfo, error){
	code, rsp,_, err:= req("job/"+job+"/lastSuccessfulBuild/api/json", []byte{})
	if err!=nil{
		panic(err)
	}
	if code!=200{
		return nil, errors.New("failed to get job details,code"+strconv.Itoa(code)+ ", "+string(rsp))
	}
	var bi BuildInfo
	err = json.Unmarshal(rsp, &bi)
	if err!=nil{
		return nil, err
	}
	return &bi, nil
}

func (jh *JHelper) build(job string,  query []byte){
	code, rsp,_, err:= req("job/"+job+"/build", query)
	if err!=nil{
		panic(err)
	}
	if code!=201{
		panic(errors.New("failed to start job details,code"+strconv.Itoa(code)+ ", "+string(rsp)))
	}
}

func (jh *JHelper) console(job string, id int, start string)(string,string, error){
	//web-rpm-build-manual/149/logText/progressiveHtml
	code, rsp,h, err:= req("job/"+job+"/"+strconv.Itoa(id)+"/logText/progressiveHtml", []byte{})
	if err!=nil{
		return "","",err
	}
	if code!=200{
		return "","",errors.New("code = "+strconv.Itoa(code))
	}
	//fmt.Println(h)
	return string(rsp),h["X-Text-Size"][0], nil
}

func req(path string, body []byte) (int, []byte,map[string][]string, error) {
	base_url :=envy.Get("URL", "URL")
	url:=base_url+path
	client := &http.Client{}
	request, err := http.NewRequest("POST", url, strings.NewReader(string(body)))
	if err != nil {
		return 0, nil,nil, err
	}
	request.Header.Add("Accept-Language", "en-us")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.SetBasicAuth(envy.Get("LOGIN", "LOGIN"), envy.Get("SECRET", "SECRET"))
	response, err := client.Do(request)
	if err != nil {
		return 0, nil,nil, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, nil,nil, err
	}

	//fmt.Println("req: ", url, "\t data: ", string(body), "type:", response.Header.Get("Content-Type"), "\n rsp: ", string(contents))
	return response.StatusCode, contents, response.Header, nil
}


