package azuredevopssdk

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func (s *Client) GetProcessTemplateList() ([]Processe, error) {
	url := fmt.Sprintf(baseURL+"%s/_apis/process/processes?api-version=5.0-preview.1", s.organization)
	log.Printf(url)
	var templateProcessList Processes
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	bytes, err := s.doRequest(req)
	log.Printf(string(bytes))
	if err != nil {
		return nil, err
	}

	json.Unmarshal(bytes, &templateProcessList)
	log.Printf("templatelist: %+v\n", templateProcessList)

	PrettyPrint(templateProcessList)

	return templateProcessList.ProcesseList, nil
}

func (s *Client) GetDefaultProcess() (*Processe, error) {

	processes, err := s.GetProcessTemplateList()
	if err != nil {
		return nil, err
	}

	for index := 0; index < len(processes); index++ {

		if processes[index].IsDefault {
			process := processes[index]
			return &process, nil
		}
	}
	return nil, nil

}

func (s *Client) GetProcessId(name string) (Processe, error) {

	var _process Processe
	processes, err := s.GetProcessTemplateList()

	if err != nil {
		return _process, err
	}

	for index := 0; index < len(processes); index++ {
		if strings.ToLower(name) == strings.ToLower(processes[index].Name) {
			process := processes[index]
			return process, nil
		}
	}
	return _process, nil

}