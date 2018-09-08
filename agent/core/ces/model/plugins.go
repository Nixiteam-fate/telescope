package model

import (
	"encoding/json"
	"github.com/huaweicloud/telescope/agent/core/logs"
	"github.com/huaweicloud/telescope/agent/core/utils"
	"io/ioutil"
	"os/exec"
)

//EachPluginConfig is the type for each plugin config
type PluginCommand struct {
	Path	string `json:"path"`
}

// PluginCmd output the plugin metric data by a plugin config
func (plugin *PluginCommand) PluginCmd() *InputMetric {

	var result InputMetric

	if !utils.IsFileExist(plugin.Path) {
		logs.GetCesLogger().Errorf("Plugin not exist: %s", plugin.Path)
		return nil
	}

	cmd := exec.Command(plugin.Path)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logs.GetCesLogger().Errorf("Plugin execute cmd StdoutPipe error: %v", err)
		return nil
	}
	defer stdout.Close()
	defer cmd.Wait()
	if err := cmd.Start(); err != nil {
		logs.GetCesLogger().Errorf("Plugin execute cmd Start error: %v", err)
		return nil
	}

	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		logs.GetCesLogger().Errorf("Plugin read all stdout error: %v", err)
		return nil
	}

	err = json.Unmarshal(opBytes, &result)
	if err != nil {
		logs.GetCesLogger().Errorf("Plugin unmarshal result error: %v", err)
		return nil
	}
	return &result
}
