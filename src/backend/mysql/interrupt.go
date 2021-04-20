package mysql

import (
	"fmt"
	"strings"

	"github.com/vj396/milton/src/types"
	"go.uber.org/zap"
)

const (
	interruptsTable = "interrupts"
)

func (c *Client) CreateInterruptRecord(i *types.Interrupt) error {
	interrupts, err := c.GetInterruptRecords(i)
	if err != nil {
		return err
	}
	for idx := range interrupts {
		if strings.EqualFold(interrupts[idx].Item, i.Item) &&
			strings.EqualFold(interrupts[idx].ChannelId, i.ChannelId) {
			return fmt.Errorf("iterrupt item %+q already in queue", i.Item)
		}
	}
	query := fmt.Sprintf("INSERT INTO %s (item, submitted_by, submitted_at, channel_id) VALUES (?, ?, ?, ?)", interruptsTable)
	stmt, err := c.Query(query, i.Item, i.SubmittedBy, i.SubmittedAt, i.ChannelId)
	defer func() {
		if err := stmt.Close(); err != nil {
			c.logger.Error("could not close insert statement op", zap.Error(err), zap.String("table", interruptsTable), zap.String("item", i.Item))
		}
	}()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetInterruptRecords(i *types.Interrupt) ([]types.Interrupt, error) {
	query := fmt.Sprintf("SELECT * FROM %s", interruptsTable)
	stmt, err := c.Query(query)
	defer func() {
		if err := stmt.Close(); err != nil {
			c.logger.Error("could not close select statement op", zap.Error(err), zap.String("table", interruptsTable))
		}
	}()
	if err != nil {
		return nil, err
	}
	var results []types.Interrupt
	for stmt.Next() {
		var interrupt types.Interrupt
		err = stmt.Scan(&interrupt.Id, &interrupt.Item, &interrupt.SubmittedBy, &interrupt.SubmittedAt, &interrupt.ChannelId)
		if err != nil {
			c.logger.Error("error scanning row.", zap.Error(err), zap.String("table", interruptsTable))
			continue
		}
		results = append(results, interrupt)
	}
	return results, nil
}

func (c *Client) DeleteInterruptRecord(i *types.Interrupt) error {
	interrupts, err := c.GetInterruptRecords(i)
	if err != nil {
		return err
	}
	itemFound := false
	for idx := range interrupts {
		if strings.EqualFold(interrupts[idx].Item, i.Item) &&
			strings.EqualFold(interrupts[idx].ChannelId, i.ChannelId) {
			i.Id = interrupts[idx].Id
			itemFound = true
			break
		}
	}
	if !itemFound {
		return fmt.Errorf("interrupt item %+q not found", i.Item)
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", interruptsTable)
	stmt, err := c.Query(query, i.Id)
	defer func() {
		if err := stmt.Close(); err != nil {
			c.logger.Error("could not close delete statement op", zap.Error(err), zap.String("table", interruptsTable), zap.String("item", i.Item))
		}
	}()
	if err != nil {
		return err
	}
	return nil
}
