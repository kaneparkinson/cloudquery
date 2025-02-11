package apprunner

import (
	"context"

	sdkTypes "github.com/cloudquery/plugin-sdk/v3/types"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/aws/aws-sdk-go-v2/service/apprunner"
	"github.com/aws/aws-sdk-go-v2/service/apprunner/types"
	"github.com/cloudquery/cloudquery/plugins/source/aws/client"
	"github.com/cloudquery/plugin-sdk/v3/schema"
	"github.com/cloudquery/plugin-sdk/v3/transformers"
)

func AutoScalingConfigurations() *schema.Table {
	tableName := "aws_apprunner_auto_scaling_configurations"
	return &schema.Table{
		Name:                tableName,
		Description:         `https://docs.aws.amazon.com/apprunner/latest/api/API_AutoScalingConfiguration.html`,
		Resolver:            fetchApprunnerAutoScalingConfigurations,
		PreResourceResolver: getAutoScalingConfiguration,
		Multiplex:           client.ServiceAccountRegionMultiplexer(tableName, "apprunner"),
		Transform:           transformers.TransformWithStruct(&types.AutoScalingConfiguration{}),
		Columns: []schema.Column{
			client.DefaultAccountIDColumn(false),
			client.DefaultRegionColumn(false),
			{
				Name:       "arn",
				Type:       arrow.BinaryTypes.String,
				Resolver:   schema.PathResolver("AutoScalingConfigurationArn"),
				PrimaryKey: true,
			},
			{
				Name:     "tags",
				Type:     sdkTypes.ExtensionTypes.JSON,
				Resolver: resolveApprunnerTags("AutoScalingConfigurationArn"),
			},
		},
	}
}

func fetchApprunnerAutoScalingConfigurations(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	var config apprunner.ListAutoScalingConfigurationsInput
	c := meta.(*client.Client)
	svc := meta.(*client.Client).Services().Apprunner
	paginator := apprunner.NewListAutoScalingConfigurationsPaginator(svc, &config)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx, func(options *apprunner.Options) {
			options.Region = c.Region
		})
		if err != nil {
			return err
		}
		res <- output.AutoScalingConfigurationSummaryList
	}
	return nil
}
func getAutoScalingConfiguration(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource) error {
	c := meta.(*client.Client)
	svc := c.Services().Apprunner
	asConfig := resource.Item.(types.AutoScalingConfigurationSummary)

	describeTaskDefinitionOutput, err := svc.DescribeAutoScalingConfiguration(ctx, &apprunner.DescribeAutoScalingConfigurationInput{AutoScalingConfigurationArn: asConfig.AutoScalingConfigurationArn}, func(options *apprunner.Options) {
		options.Region = c.Region
	})
	if err != nil {
		return err
	}

	resource.Item = describeTaskDefinitionOutput.AutoScalingConfiguration
	return nil
}
