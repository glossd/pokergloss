package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/glossd/pokergloss/survival/bot-squad/conf"
	"github.com/glossd/pokergloss/survival/bot-squad/domain"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func GetTable(config conf.Config) (*domain.Table, error) {
	url := fmt.Sprintf("%s://%s/api/table/tables/%s", config.Table.Scheme, config.Table.Host, config.TableID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %s", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", config.Token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch the table: %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read the table body: %s", err)
	}

	var t domain.Table
	err = json.Unmarshal(body, &t)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal the table: %s", err)
	}
	return &t, nil
}

func MakeAction(config conf.Config, position int, a domain.Action, t *domain.Table) {
	url := fmt.Sprintf("%s://%s/api/table/tables/%s/seats/%d/actions/%s", config.Table.Scheme, config.Table.Host, config.TableID, position, a.Type)
	body := bytes.NewBuffer([]byte(fmt.Sprintf(`{"chips":%d}`, a.Chips)))
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		fmt.Errorf("Failed to create request: %s", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", config.Token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Failed to do make action request: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		msg, _ := ioutil.ReadAll(resp.Body)
		log.Errorf("Make action failed with status %s: %s", resp.Status, msg)
	}
}

func SitBack(config conf.Config, position int) {
	url := fmt.Sprintf("%s://%s/api/table/tables/%s/seats/%d/sit-back", config.Table.Scheme, config.Table.Host, config.TableID, position)
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		fmt.Errorf("Failed to create request: %s", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", config.Token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Failed to do sit back request: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		msg, _ := ioutil.ReadAll(resp.Body)
		log.Errorf("Sit back failed with status %s: %s", resp.Status, msg)
	}
}
