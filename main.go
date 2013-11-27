
package main

import (
    "fmt"
    "os"
    "flag"
    "encoding/xml"
    "encoding/csv"
    "bytes"
    "text/template"
    "text/tabwriter"
)


type TemplateConfig struct {
    Name string `xml:"name,attr"`
    Lines []string `xml:"Lines>Line"`
}

type FileConfig struct {
    Name string `xml:"name,attr"`
    TemplateName string `xml:"templateName,attr"`
}

type Config struct {
    TemplateConfig []TemplateConfig
    FileConfig []FileConfig
}

func PadBefore(s string,i int) string {
    if s=="" {
        return ""
    }

    for j:=0; j<i; j=j+1 {
        s=" " + s
    }
    return s
}

func makeTemplates(t *[]TemplateConfig) (*template.Template,error) {
    var textTemplates *template.Template
    var err error
    templateFuncs := template.FuncMap{"padBefore":PadBefore}
    for i,templateConfig := range *t {
        t:=""
        for i,l:=range templateConfig.Lines {
            if i>0 {t=t+"\t"}
            t=t+l
        }
        if i==0 {
            textTemplates,err=template.New(templateConfig.Name).Funcs(templateFuncs).Parse(t)
        } else {
            textTemplates,err=textTemplates.New(templateConfig.Name).Funcs(templateFuncs).Parse(t)
        }
        if err!=nil {break}
    }
    return textTemplates,err
}

func makeTemplateContext(r *[]string) (map[string]string) {
    context:=make(map[string]string,len(*r))
    for i,v:=range *r {
        // i starts at 0, but we want to start with A which is 65. 
        context["c"+string(i+65)]=v
    }
    return context
}


func processHeader(fc FileConfig, t *template.Template) () {
        // in file
    iFile,err := os.Open(fc.Name)
    if err!=nil {panic("Can't open " + fc.Name)}
    defer iFile.Close()

    r := csv.NewReader(iFile)
    r.Comma = '\t'
    r.TrailingComma = true
    record,err:=r.Read()
    fmt.Println("File: ",fc.Name)
    w:= new(tabwriter.Writer)
    w.Init(os.Stdout, 0, 8, 0, '\t', tabwriter.Debug)
    t.ExecuteTemplate(w,fc.TemplateName,makeTemplateContext(&record))
    w.Flush()
    fmt.Println("")
}

func processFile(fc FileConfig, t *template.Template) () {
    // in file
    iFile,err := os.Open(fc.Name)
    if err!=nil {panic("Can't open " + fc.Name)}
    defer iFile.Close()

    // out file
    oFile,err := os.Create(fc.Name + "_out.txt")
    if err!=nil {panic("Can't open/write to " + fc.Name + "_out.txt")}
    defer oFile.Close()

    r := csv.NewReader(iFile)
    r.Comma = '\t'
    r.TrailingComma = true
    record,err:=r.Read()

    for i:=0;err==nil;record,err=r.Read() {
        if i>0 {
            oFile.Write([]byte("\n"))
        }
        t.ExecuteTemplate(oFile,fc.TemplateName,makeTemplateContext(&record))
        i=i+1
    }
    fmt.Println(err)
    oFile.Sync()
}


func main() {
    var (
        config Config
    )

    configFileName := flag.String("config", "","Config file")
    headerReport := flag.Bool("headerReport", false,"Produce a report of header mappings")
    flag.Parse()

    configFile, err := os.Open(*configFileName)
    if err != nil {
        fmt.Println("Error opening config:", err)
        return
    }
    defer configFile.Close()
    configBuf := new(bytes.Buffer)
    configBuf.ReadFrom(configFile)
    xml.Unmarshal(configBuf.Bytes(),&config)

    // Parse templates
    textTemplates,err := makeTemplates(&config.TemplateConfig)
    if(err!=nil){panic(err)}

    // Process each input file config
    for _,fileConfig := range config.FileConfig {
        if *headerReport {
            processHeader(fileConfig,textTemplates)
            continue
        }
        processFile(fileConfig,textTemplates)


    }
}
