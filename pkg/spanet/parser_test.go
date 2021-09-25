package spanet_test

import (
	"strings"
	"testing"

	"git.fremnet.net/spanet/pkg/spanet"
	"github.com/davecgh/go-spew/spew"
)

const testRF = `RF:
,R2,18,250,51,70,4,13,50,55,19,6,2020,376,9999,1,0,490,207,34,6000,602,23,20,0,0,0,0,44,35,45,:
,R3,32,1,4,4,4,SW V5 17 05 31,SV3,18480001,20000826,1,0,0,0,0,0,NA,7,0,470,Filtering,4,0,7,7,0,0,:
,R4,NORM,0,0,0,1,0,3547,4,20,4500,7413,567,1686,0,8388608,0,0,5,0,98,0,10084,4,80,100,0,0,4,:
,R5,0,1,0,1,0,0,0,0,0,0,1,0,1,0,376,0,3,4,0,0,0,0,0,1,2,6,:
,R6,1,5,30,2,5,8,1,360,1,0,3584,5120,127,128,5632,5632,2304,1792,0,30,0,0,0,0,2,3,0,:
,R7,2304,0,1,1,1,0,1,0,0,0,253,191,253,240,483,125,77,1,0,0,0,23,200,1,0,1,31,32,35,100,5,:
,R9,F1,255,0,0,0,0,0,0,0,0,0,0,:
,RA,F2,0,0,0,0,0,0,255,0,0,0,0,:
,RB,F3,0,0,0,0,0,0,0,0,0,0,0,:
,RC,0,1,1,0,0,0,0,0,0,2,0,0,1,0,:
,RE,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,-4,13,30,8,5,1,0,0,0,0,0,:*
,RG,1,1,1,1,1,1,1-1-014,1-1-01,1-1-01,0-,0-,0,:*`

func TestThing(t *testing.T) {
	net, err := spanet.Parse(strings.NewReader(testRF))
	if err != nil {
		panic(err)
	}
	spew.Dump(net)
}
