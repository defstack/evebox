/* Copyright (c) 2016 Jason Ish
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 *
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED ``AS IS'' AND ANY EXPRESS OR IMPLIED
 * WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT,
 * INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
 * STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING
 * IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */

package eve

import (
	"github.com/jasonish/evebox/geoip"
	"github.com/jasonish/evebox/log"
	"net"
)

var RFC1918_Netstrings = []string{
	"10.0.0.0/8",
	"127.16.0.0/12",
	"192.168.0.0/16",
}

var RFC1918_IPNets []*net.IPNet

func IsRFC1918(addr string) bool {
	ip := net.ParseIP(addr)
	for _, ipnet := range RFC1918_IPNets {
		if ipnet.Contains(ip) {
			return true
		}
	}
	return false
}

func init() {
	for _, network := range RFC1918_Netstrings {
		_, ipnet, err := net.ParseCIDR(network)
		if err == nil {
			RFC1918_IPNets = append(RFC1918_IPNets, ipnet)
		}
	}
}

type GeoipFilter struct {
	db *geoip.GeoIpDb
}

func NewGeoipFilter(db *geoip.GeoIpDb) *GeoipFilter {
	return &GeoipFilter{
		db: db,
	}
}

func (f *GeoipFilter) AddGeoIP(event RawEveEvent) {

	if f.db == nil {
		return
	}

	srcip, ok := event["src_ip"].(string)
	if ok && !IsRFC1918(srcip) {
		gip, err := f.db.LookupString(srcip)
		if err != nil {
			log.Debug("Failed to lookup geoip for %s", srcip)
		}

		// Need at least a continent code.
		if gip.ContinentCode != "" {
			event["geoip"] = gip
		}
	}
	if event["geoip"] == nil {
		destip, ok := event["dest_ip"].(string)
		if ok && !IsRFC1918(destip) {
			gip, err := f.db.LookupString(destip)
			if err != nil {
				log.Debug("Failed to lookup geoip for %s", destip)
			}
			// Need at least a continent code.
			if gip.ContinentCode != "" {
				event["geoip"] = gip
			}
		}
	}

}
