package receiver

import (
	"fmt"
	"io"

	collectorpb "github.com/Percona-Lab/qan-api/api/collector"
	"github.com/Percona-Lab/qan-api/models"
)

type Service struct {
	qcm models.QueryClass
}

func NewService(qcm models.QueryClass) *Service {
	return &Service{qcm}
}

func (s *Service) DataInterchange(stream collectorpb.Agent_DataInterchangeServer) error {
	fmt.Println("Start...")
	for {
		agentMsg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("recved from agent: %+v", agentMsg)
		}
		err = s.qcm.Save(agentMsg)
		if err != nil {
			return fmt.Errorf("save error: %s", err.Error())
		}
		fmt.Printf("Rcvd and saved %v QC\n", len(agentMsg.QueryClass))
	}
}
