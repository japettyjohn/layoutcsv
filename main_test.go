package main

import (
    "testing"
    "encoding/xml"
)

func TestMain(t *testing.T) {

    sampleConfig:=`<config>
    <TemplateConfig name="template 1">
        <Lines>
            <Line header="ID">{{.cA}}</Line>
            <Line>{{.cB}}</Line>
            <Line>{{.cC}}</Line>
            <Line>{{.cD}}</Line>
            <Line>{{.cE}}</Line>
            <Line>{{.cF}}</Line>
            <Line>{{.cG}}</Line>
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
        if v.Value!=("{{.c"+string(i+65) + "}}") {
            t.Log("Test was: ","col"+string(i+65))
            t.Fatal("Order is wrong: ",v," index: ",i)
        }
    }
    if config.TemplateConfig[0].Lines[0].Header!="ID" {
        t.Fatal("Wrong header: ",config.TemplateConfig[0].Lines[0].Header)
    }
    


    t.Log("Pass: ", config)
}