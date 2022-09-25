// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package inputSelfDescribing

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/buz/pkg/annotator"
	"github.com/silverton-io/buz/pkg/params"
	"github.com/silverton-io/buz/pkg/privacy"
	"github.com/silverton-io/buz/pkg/response"
)

func Handler(h params.Handler) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		if c.ContentType() == "application/json" {
			envelopes := BuildEnvelopesFromRequest(c, h.Config, h.CollectorMeta)
			annotatedEnvelopes := annotator.Annotate(envelopes, h.Registry)
			anonymizedEnvelopes := privacy.AnonymizeEnvelopes(annotatedEnvelopes, h.Config.Privacy)
			err := h.Manifold.Distribute(anonymizedEnvelopes, h.ProtocolStats)
			if err != nil {
				c.Header("Retry-After", response.RETRY_AFTER_60)
				c.JSON(http.StatusServiceUnavailable, response.ManifoldDistributionError)
			} else {
				c.JSON(200, response.Ok)
			}

		} else {
			c.JSON(http.StatusBadRequest, response.InvalidContentType)
		}
	}
	return gin.HandlerFunc(fn)
}