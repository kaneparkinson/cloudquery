package directconnect

import (
	"context"

	sdkTypes "github.com/cloudquery/plugin-sdk/v3/types"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/aws/aws-sdk-go-v2/service/directconnect"
	"github.com/aws/aws-sdk-go-v2/service/directconnect/types"
	"github.com/cloudquery/cloudquery/plugins/source/aws/client"
	"github.com/cloudquery/plugin-sdk/v3/schema"
	"github.com/cloudquery/plugin-sdk/v3/transformers"
)

func Lags() *schema.Table {
	tableName := "aws_directconnect_lags"
	return &schema.Table{
		Name:        tableName,
		Description: `https://docs.aws.amazon.com/directconnect/latest/APIReference/API_Lag.html`,
		Resolver:    fetchDirectconnectLags,
		Multiplex:   client.ServiceAccountRegionMultiplexer(tableName, "directconnect"),
		Transform:   transformers.TransformWithStruct(&types.Lag{}),
		Columns: []schema.Column{
			client.DefaultAccountIDColumn(false),
			client.DefaultRegionColumn(false),
			{
				Name:       "arn",
				Type:       arrow.BinaryTypes.String,
				Resolver:   resolveLagARN(),
				PrimaryKey: true,
			},
			{
				Name:     "id",
				Type:     arrow.BinaryTypes.String,
				Resolver: schema.PathResolver("LagId"),
			},
			{
				Name:     "tags",
				Type:     sdkTypes.ExtensionTypes.JSON,
				Resolver: client.ResolveTags,
			},
		},
	}
}

func fetchDirectconnectLags(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	var config directconnect.DescribeLagsInput
	c := meta.(*client.Client)
	svc := c.Services().Directconnect
	output, err := svc.DescribeLags(ctx, &config, func(options *directconnect.Options) {
		options.Region = c.Region
	})
	if err != nil {
		return err
	}
	res <- output.Lags
	return nil
}

func resolveLagARN() schema.ColumnResolver {
	return client.ResolveARN(client.DirectConnectService, func(resource *schema.Resource) ([]string, error) {
		return []string{"dxlag", *resource.Item.(types.Lag).LagId}, nil
	})
}
