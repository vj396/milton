package mysql

import (
	"fmt"

	"github.com/vj396/milton/src/types"
	"go.uber.org/zap"
)

const (
	documentsTable = "documents"
)

func (c *Client) CreateDocmentsRecord(d *types.Documents) error {
	documents, err := c.GetDocumentsRecords(d)
	if err != nil {
		return err
	}
	for idx := range documents {
		if documents[idx].Title == d.Title {
			return fmt.Errorf("document entry already present: %+q, %+q", d.Title, d.Link)
		}
	}

	query := fmt.Sprintf("INSERT INTO %s (title, link) VALUES (?, ?)", documentsTable)
	stmt, err := c.Query(query, d.Title, d.Link)
	defer func() {
		if err := stmt.Close(); err != nil {
			c.logger.Error("could not close insert statement op", zap.Error(err), zap.String("table", documentsTable), zap.String("document-title", d.Title))
		}
	}()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetDocumentsRecords(d *types.Documents) ([]types.Documents, error) {
	query := fmt.Sprintf("SELECT * FROM %s", documentsTable)
	stmt, err := c.Query(query)
	defer func() {
		if err := stmt.Close(); err != nil {
			c.logger.Error("could not close select statement op", zap.Error(err), zap.String("table", documentsTable))
		}
	}()
	if err != nil {
		return nil, err
	}
	var results []types.Documents
	for stmt.Next() {
		var doc types.Documents
		err = stmt.Scan(&doc.Id, &doc.Title, &doc.Link)
		if err != nil {
			c.logger.Error("error scanning row.", zap.Error(err), zap.String("table", documentsTable))
			continue
		}
		results = append(results, doc)
	}
	return results, nil
}

func (c *Client) DeleteDocumentsRecord(d *types.Documents) error {
	documents, err := c.GetDocumentsRecords(d)
	if err != nil {
		return err
	}
	documentFound := false
	for idx := range documents {
		if documents[idx].Title == d.Title {
			d.Id = documents[idx].Id
			documentFound = true
			break
		}
	}
	if !documentFound {
		return fmt.Errorf("no document titled %+q found", d.Title)
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", documentsTable)
	stmt, err := c.Query(query, d.Id)
	defer func() {
		if err := stmt.Close(); err != nil {
			c.logger.Error("could not close delete statement op", zap.Error(err), zap.String("table", documentsTable), zap.String("document-title", d.Title))
		}
	}()
	if err != nil {
		return err
	}
	return nil
}
