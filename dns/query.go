package dns

import (
	"github.com/jaeles-project/jaeles/libs"
	"github.com/jaeles-project/jaeles/utils"
	"github.com/lixiangzhong/dnsutil"
	"github.com/thoas/go-funk"
)

var recordMap = map[string]uint16{
	"A":     1,
	"AAAA":  28,
	"NS":    2,
	"CNAME": 5,
	"SOA":   6,
	"PTR":   12,
	"MX":    15,
	"TXT":   16,
}

var CommonResolvers = []string{
	"1.1.1.1", // Cloudflare
	"8.8.8.8", // Google
	"8.8.4.4", // Google
}

func QueryDNS(dnsRecord *libs.Dns, options libs.Options) {
	resolver := options.Resolver
	if resolver == "" {
		index := funk.RandomInt(0, len(CommonResolvers))
		resolver = CommonResolvers[index]
	}
	domain := dnsRecord.Domain
	queryType := dnsRecord.RecordType

	var dig dnsutil.Dig
	dig.Retry = options.Retry
	dig.SetDNS(resolver)
	utils.InforF("[resolved] %v -- %v", domain, queryType)

	if queryType == "ANY" || queryType == "" {
		for k, v := range recordMap {
			var dnsResult libs.DnsResult
			msg, err := dig.GetMsg(v, domain)
			if err != nil {
				utils.DebugF("err to resolved: %v -- %v", domain, err)
				return
			}
			dnsResult.Data = msg.String()
			dnsResult.RecordType = k
			dnsRecord.Results = append(dnsRecord.Results, dnsResult)
		}
	} else {
		var dnsResult libs.DnsResult
		msg, err := dig.GetMsg(recordMap[queryType], domain)
		if err != nil {
			utils.DebugF("err to resolved: %v -- %v", domain, err)
			return
		}
		dnsResult.Data = msg.String()
		dnsResult.RecordType = queryType
		dnsRecord.Results = append(dnsRecord.Results, dnsResult)
	}

	return
}
