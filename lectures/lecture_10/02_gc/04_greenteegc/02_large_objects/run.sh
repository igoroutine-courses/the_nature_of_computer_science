echo "old: "
go test -bench=.



echo "greenteagc: "
(GOEXPERIMENT=greenteagc go test -bench=.)
