package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ipfs/testground/pkg/api"
	"github.com/ipfs/testground/pkg/daemon/client"
	"github.com/ipfs/testground/pkg/tgwriter"
	"go.uber.org/zap"
)

var TermExplanation = "a term is any of: <testplan> or <testplan>/<testcase>"

func (srv *Server) describeHandler(w http.ResponseWriter, r *http.Request, log *zap.SugaredLogger) {
	log.Debugw("handle request", "command", "describe")
	defer log.Debugw("request handled", "command", "describe")

	tgw := tgwriter.New(w, log)

	var req client.DescribeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		tgw.WriteError("cannot json decode request body", "err", err)
		return
	}

	term := req.Term

	engine, err := GetEngine()
	if err != nil {
		tgw.WriteError("get engine error", "err", err)
		return
	}

	var pl, tc string
	switch splt := strings.Split(term, "/"); len(splt) {
	case 2:
		pl, tc = splt[0], splt[1]
	case 1:
		pl = splt[0]
	default:
		tgw.WriteError("unrecognized format for term", "explanation", TermExplanation)
		return
	}

	plan := engine.TestCensus().PlanByName(pl)
	if plan == nil {
		tgw.WriteError(fmt.Sprintf("plan not found, name: %s ; term: %s", pl, term))
		return
	}

	var cases []*api.TestCase
	if tc == "" {
		cases = plan.TestCases
	} else if _, tcbn, ok := plan.TestCaseByName(tc); ok {
		cases = []*api.TestCase{tcbn}
	} else {
		tgw.WriteError(fmt.Sprintf("test case not found: %s", tc))
		return
	}

	plan.Describe(tgw)

	header := `TESTCASES:
----------
----------
`

	_, err = tgw.Write([]byte(header))
	if err != nil {
		tgw.WriteError("header write error", "err", err)
		return
	}

	for _, tc := range cases {
		tc.Describe(tgw)
	}

	tgw.WriteResult(struct{}{})
}
