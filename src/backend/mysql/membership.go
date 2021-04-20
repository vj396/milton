package mysql

import (
	"fmt"
	"time"

	"github.com/vj396/milton/src/types"
	"go.uber.org/zap"
)

const (
	membershipTable = "membership"
)

var (
	defaultMembershipSLA = time.Duration(24 * time.Hour)
)

func (c *Client) CreateMembershipRecord(m *types.Membership) error {
	memberships, err := c.GetMemebershipRecords(m)
	if err != nil {
		return err
	}
	for idx := range memberships {
		if memberships[idx].Id == m.Id {
			return fmt.Errorf("bot already a member of channel: %+q", m.ChannelName)
		}
	}
	if m.SLA == 0 {
		m.SLA = uint(defaultMembershipSLA)
	}
	query := fmt.Sprintf("INSERT INTO %s (id, sla, channel_name) VALUES (?, ?, ?)", membershipTable)
	stmt, err := c.Query(query, m.Id, m.SLA, m.ChannelName)
	defer func() {
		if err := stmt.Close(); err != nil {
			c.logger.Error("could not close insert statement op", zap.Error(err), zap.String("table", membershipTable), zap.String("channel-name", m.ChannelName))
		}
	}()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetMemebershipRecords(m *types.Membership) ([]types.Membership, error) {
	query := fmt.Sprintf("SELECT * FROM %s", membershipTable)
	stmt, err := c.Query(query)
	defer func() {
		if err := stmt.Close(); err != nil {
			c.logger.Error("could not close select statement op", zap.Error(err), zap.String("table", membershipTable))
		}
	}()
	if err != nil {
		return nil, err
	}
	var results []types.Membership
	for stmt.Next() {
		var membership types.Membership
		err = stmt.Scan(&membership.Id, &membership.SLA, &membership.ChannelName)
		if err != nil {
			c.logger.Error("error scanning row.", zap.Error(err), zap.String("table", membershipTable))
			continue
		}
		results = append(results, membership)
	}
	return results, nil
}

func (c *Client) DeleteMembershipRecord(m *types.Membership) error {
	memberships, err := c.GetMemebershipRecords(m)
	if err != nil {
		return err
	}
	channelFound := false
	for idx := range memberships {
		if memberships[idx].Id == m.Id {
			channelFound = true
			break
		}
	}
	if !channelFound {
		return fmt.Errorf("bot not a member of channel, %+v", m.ChannelName)
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", membershipTable)
	stmt, err := c.Query(query, m.Id)
	defer func() {
		if err := stmt.Close(); err != nil {
			c.logger.Error("could not close delete statement op", zap.Error(err), zap.String("table", membershipTable), zap.String("channel-name", m.ChannelName))
		}
	}()
	if err != nil {
		return err
	}
	return nil
}
