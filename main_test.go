package main

import (
    "testing"
    "encoding/xml"
    "strconv"
)

func TestMain(t *testing.T) {

    sampleConfig:=`<config>
    <TemplateConfig name="template 1">
        <Lines>
            <Line>{{.c1}}</Line>
            <Line>{{.c2}}</Line>
            <Line>{{.c3}}</Line>
            <Line>{{.c4}}</Line>
            <Line>{{.c5}}</Line>
            <Line>{{.c6}}</Line>
            <Line>{{.c7}}</Line>
        </Lines>
    </TemplateConfig>
    <FileConfig name="aFileName.xls" templateName="template 1" />
</config>`

    var (
        config Config
    )

    err:=xml.Unmarshal([]byte(sampleConfig),&config)
    if err!=nil {
        t.Fatal("Error parsing",err)
    }

    if len(config.FileConfig) != 1 {
        t.Fatal("FileConfig not parsed: ",config)
    }
    if !t.Failed() && config.FileConfig[0].Name != "aFileName.xls" {
        t.Fatal("FileConfig name!=aFileName.xls")
    }

    for i,v:= range config.TemplateConfig[0].Lines {
        if v!=("{{.c"+strconv.Itoa(i+1) + "}}") {
            t.Log("Test was: ","col"+strconv.Itoa(i+1))
            t.Fatal("Order is wrong: ",v," index: ",i)
        }
    }
    


    t.Log("Pass: ", config)
}