package sharepoint

import (
	"testing"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/alcionai/corso/src/pkg/backup/details"
)

type SharePointInfoSuite struct {
	suite.Suite
}

func TestSharePointInfoSuite(t *testing.T) {
	suite.Run(t, new(SharePointInfoSuite))
}

func (suite *SharePointInfoSuite) TestSharePointInfo() {
	tests := []struct {
		name      string
		listAndRP func() (models.Listable, *details.SharepointInfo)
	}{
		{
			name: "Empty List",
			listAndRP: func() (models.Listable, *details.SharepointInfo) {
				i := &details.SharepointInfo{ItemType: details.SharepointItem}
				return models.NewList(), i
			},
		}, {
			name: "Only Name",
			listAndRP: func() (models.Listable, *details.SharepointInfo) {
				aTitle := "Whole List"
				listing := models.NewList()
				listing.SetDisplayName(&aTitle)
				i := &details.SharepointInfo{
					ItemType: details.SharepointItem,
					ItemName: aTitle,
				}
				return listing, i
			},
		},
	}
	for _, test := range tests {
		suite.T().Run(test.name, func(t *testing.T) {
			list, expected := test.listAndRP()
			info := sharepointListInfo(list)
			assert.Equal(t, expected.ItemType, info.ItemType)
			assert.Equal(t, expected.ItemName, info.ItemName)
			assert.Equal(t, expected.WebURL, info.WebURL)
		})
	}
}
