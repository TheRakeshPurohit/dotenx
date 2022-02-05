package jobService

import (
	"fmt"

	"github.com/utopiops/automated-ops/runner/config"
	"github.com/utopiops/automated-ops/runner/executors"
	"github.com/utopiops/automated-ops/runner/models"
	"github.com/utopiops/automated-ops/runner/shared"
)

func (manager *JobManager) HandleJob(job models.Job, logHelper shared.LogHelper) {
	executor := executors.NewExecutor()
	var name, taskType, service string
	var body map[string]interface{}
	if _, ok := job.Data["name"]; ok {
		name = job.Data["name"].(string)
	}
	if _, ok := job.Data["type"]; ok {
		taskType = job.Data["type"].(string)
	}
	if _, ok := job.Data["body"]; ok {
		body = job.Data["body"].(map[string]interface{})
	}
	if _, ok := job.Data["serviceAccount"]; ok {
		service = job.Data["serviceAccount"].(string)
	}
	taskDetails := models.TaskDetails{
		Name:           name,
		Id:             job.Id,
		Type:           taskType,
		Body:           body,
		ServiceAccount: service,
	}
	result := executor.Execute(&taskDetails)
	fmt.Println(result)
	resultDto := models.TaskStatus{
		ReturnValue: result.Log,
		Result:      models.Status(result.Status),
		Toekn:       config.Configs.Queue.Token,
	}
	var err error
	var id string
	if result.Error == nil {
		resultDto.Result = models.StatusCompleted
		id, err = manager.LogHelper.Log("log: "+result.Log, true, result.Id)
	} else {
		id, err = manager.LogHelper.Log("error: "+result.Error.Error()+", log: "+result.Log, true, result.Id)
		resultDto.Result = models.StatusFailed
	}
	if err != nil {
		fmt.Printf("error in setting job log: %s\n", err.Error())
	} else {
		fmt.Println("jobId: " + id)
	}
	manager.SendResult(job.Id, resultDto)
}
