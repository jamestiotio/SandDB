package db

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"sanddb/read_write"
)

func (h *Handler) HandleCreateTable(c *fiber.Ctx) error {
	var (
		reqBody read_write.CreateRequest
	)

	//TODO: check if table name already exists
	err := c.BodyParser(&reqBody)
	if err != nil {
		_ = c.SendStatus(http.StatusInternalServerError)
		return err
	}
	partitions := make([]Partition, 0)
	table := Table{
		TableName:          reqBody.TableName,
		PartitionKeyNames:  reqBody.PartitionKeyNames,
		ClusteringKeyNames: reqBody.ClusteringKeyNames,
		Partitions:         partitions,
	}

	filename := fmt.Sprintf("%d/%d.json", h.Node.Id, h.Node.Id)
	err = PersistTable(filename, table)
	if err != nil {
		_ = c.SendStatus(http.StatusInternalServerError)
		return err
	}
	//TODO: reply to the coordinator that node manages to create table
	responseMsg := &read_write.PeerMessage{
		Type:     read_write.CREATE_ACK,
		Content:  "1",
		SourceID: h.Node.Id,
	}
	resp, err := json.Marshal(responseMsg)
	if err != nil {
		_ = c.SendStatus(http.StatusInternalServerError)
		return err
	}
	_ = c.Send(resp)
	fmt.Printf("Finished creating Table: %s", table.TableName)
	return nil
}
