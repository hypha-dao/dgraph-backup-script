package gql

import (
	"context"
	"fmt"
	"time"

	"github.com/machinebox/graphql"
)

const ResponseCode_Success = "Success"

type ExportArgs struct {
	Format      string
	Destination string
	AccessKey   string
	SecretKey   string
}

func (m *ExportArgs) String() string {
	return fmt.Sprintf(
		`
			ExportArgs {
				Format: %v
				Destination: %v
				AccessKey: %v
				SecretKey: %v
			}	
		`,
		m.Format,
		m.Destination,
		m.AccessKey,
		m.SecretKey,
	)
}

type TaskStatus struct {
	Status      string
	LastUpdated time.Time
	Kind        string
}

type Admin struct {
	client   *graphql.Client
	Endpoint string
}

func NewAdmin(endpoint string) *Admin {
	return &Admin{
		client:   graphql.NewClient(endpoint),
		Endpoint: endpoint,
	}
}

func (m *Admin) GetTaskStatus(taskId string) (*TaskStatus, error) {
	req := graphql.NewRequest(`
		query($taskId: String!) {
			task(input: {id: $taskId}) {
        status
        lastUpdated
        kind
    	}
		}
	`)
	req.Var("taskId", taskId)

	var response interface{}
	err := m.client.Run(context.Background(), req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed getting task status, error: %v", err)
	}
	fmt.Println("Response: ", response)
	task := response.(map[string]interface{})["task"].(map[string]interface{})

	return &TaskStatus{
		Status:      task["status"].(string),
		LastUpdated: task["lastUpdated"].(time.Time),
		Kind:        task["kind"].(string),
	}, nil
}

func (m *Admin) Export(args *ExportArgs) error {

	if args.Format == "" {
		args.Format = "rdf"
	}

	req := graphql.NewRequest(`
		mutation($format: String!, $destination: String!, $accessKey: String!, $secretKey: String!) {
			export(input: {
				format: $format
				destination: $destination
				accessKey: $accessKey
				secretKey: $secretKey
			}) {
				response {
					message
					code
				}
			}
		}
	`)
	req.Var("format", args.Format)
	req.Var("destination", args.Destination)
	req.Var("accessKey", args.AccessKey)
	req.Var("secretKey", args.SecretKey)

	var response interface{}
	err := m.client.Run(context.Background(), req, &response)
	if err != nil {
		return fmt.Errorf("failed exporting DB, error: %v", err)
	}

	fmt.Println("Response: ", response)
	resp := response.(map[string]interface{})["export"].(map[string]interface{})["response"].(map[string]interface{})
	code := resp["code"].(string)
	if code != ResponseCode_Success {
		return fmt.Errorf("failed exporting DB, error: %v", resp["message"].(string))
	}
	health, err := m.Health()
	fmt.Println("Health: ", health)
	return nil
}

func (m *Admin) Health() (string, error) {
	req := graphql.NewRequest(`
		{
			health{
				instance
				status
				ongoing
				indexing
			}
		}
	`)
	var response interface{}
	err := m.client.Run(context.Background(), req, &response)
	if err != nil {
		return "", fmt.Errorf("failed getting health state, error: %v", err)
	}
	return fmt.Sprintf("%v", response.(map[string]interface{})["health"]), nil
}

func (m *Admin) String() string {
	return fmt.Sprintf(
		`
			Admin:{
				Endpoint: %v
			}
		`,
		m.Endpoint,
	)
}
