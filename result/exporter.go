package result

import(
"os"
"bufio"
"encoding/gob"
// "encoding/json"
)



func SaveToFile(result *Result, fileName string) {
    encodeFile, err := os.Create(fileName)
    if err != nil {
        panic(err)
    }
    bufW := bufio.NewWriter(encodeFile)
    // Since this is a binary format large parts of it will be unreadable
    // encoder := json.NewEncoder(bufW)
       encoder := gob.NewEncoder(bufW)


    // Write to the file
    if err := encoder.Encode(*result); err != nil {
        panic(err)
    }
    encodeFile.Close()
}

