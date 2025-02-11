package network

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"github.com/cloudquery/cloudquery/plugins/source/azure/client"
	"github.com/cloudquery/plugin-sdk/v3/schema"
	"github.com/cloudquery/plugin-sdk/v3/transformers"
)

func watcherFlowLogs() *schema.Table {
	return &schema.Table{
		Name:        "azure_network_watcher_flow_logs",
		Resolver:    fetchWatcherFlowLogs,
		Description: "https://learn.microsoft.com/en-us/rest/api/network-watcher/flow-logs/list?tabs=HTTP#definitions",
		Transform:   transformers.TransformWithStruct(&armnetwork.FlowLog{}, transformers.WithPrimaryKeys("ID")),
		Columns:     schema.ColumnList{client.SubscriptionID},
	}
}

func fetchWatcherFlowLogs(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	p := parent.Item.(*armnetwork.Watcher)
	cl := meta.(*client.Client)
	svc, err := armnetwork.NewFlowLogsClient(cl.SubscriptionId, cl.Creds, cl.Options)
	if err != nil {
		return err
	}
	group, err := client.ParseResourceGroup(*p.ID)
	if err != nil {
		return err
	}
	pager := svc.NewListPager(group, *p.Name, nil)
	for pager.More() {
		p, err := pager.NextPage(ctx)
		if err != nil {
			return err
		}
		res <- p.Value
	}
	return nil
}
