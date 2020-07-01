package cmd

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetJobDetails(t *testing.T) {
	jh:= CreateJHelper()
	fmt.Println("jd ", jh.jobDetails)
}


func TestGetLastSuccessfulBuildDuration(t *testing.T) {
	jh:= CreateJHelper()
	rsp, err := jh.getLastSuccessfulBuildInfo("config-deploy-manual")
	assert.NoError(t, err)
	fmt.Println("jd ", rsp)
}


func TestGetJobInfo(t *testing.T) {
	jh:= CreateJHelper()
	ji := jh.getJobInfo("core-change-zone")
	fmt.Printf("ji %+v", ji.getParameterDefinitions())
}


