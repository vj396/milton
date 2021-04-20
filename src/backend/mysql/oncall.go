package mysql

import (
	"fmt"
	"strings"

	"github.com/vj396/milton/src/types"
	"go.uber.org/zap"
)

const (
	oncallTable = "oncall"
)

func (c *Client) CreateOncallRecord(o *types.Oncall) error {
	oncalls, err := c.GetOncallRecords(o)
	if err != nil {
		return err
	}
	for idx := range oncalls {
		if (oncalls[idx].Id == o.Id) && strings.EqualFold(oncalls[idx].Type, o.Type) {
			return fmt.Errorf("oncall team with id %+q already being tracked", o.Id)
		}
	}
	query := fmt.Sprintf("INSERT INTO %s (identifier, team_name, team_type) VALUES (?, ?, ?)", oncallTable)
	stmt, err := c.Query(query, o.Id, o.Name, o.Type)
	defer func() {
		if err := stmt.Close(); err != nil {
			c.logger.Error("could not close insert statement op", zap.Error(err), zap.String("table", oncallTable), zap.String("id", o.Id), zap.String("name", o.Name), zap.String("backend", o.Type))
		}
	}()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetOncallRecords(o *types.Oncall) ([]types.Oncall, error) {
	query := fmt.Sprintf("SELECT * FROM %s", oncallTable)
	stmt, err := c.Query(query)
	defer func() {
		if err := stmt.Close(); err != nil {
			c.logger.Error("could not close select statement op", zap.Error(err), zap.String("table", oncallTable))
		}
	}()
	if err != nil {
		return nil, err
	}
	var results []types.Oncall
	for stmt.Next() {
		var oncall types.Oncall
		err = stmt.Scan(&oncall.Id, &oncall.Name, &oncall.Type)
		if err != nil {
			c.logger.Error("error scanning row.", zap.Error(err), zap.String("table", oncallTable))
			continue
		}
		results = append(results, oncall)
	}
	return results, nil
}

func (c *Client) DeleteOncallRecord(o *types.Oncall) error {
	oncalls, err := c.GetOncallRecords(o)
	if err != nil {
		return err
	}
	oncallFound := false
	for idx := range oncalls {
		if (oncalls[idx].Id == o.Id) && strings.EqualFold(oncalls[idx].Type, o.Type) {
			oncallFound = true
			break
		}
	}
	if !oncallFound {
		return fmt.Errorf("no oncall team with ID %+q in %+q not found", o.Id, o.Type)
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE identifier = ? and team_type = ?", oncallTable)
	stmt, err := c.Query(query, o.Id, o.Type)
	defer func() {
		if err := stmt.Close(); err != nil {
			c.logger.Error("could not close delete statement op", zap.Error(err), zap.String("table", oncallTable), zap.String("oncall-team", o.Name), zap.String("backend", o.Type))
		}
	}()
	if err != nil {
		return err
	}
	return nil
}
