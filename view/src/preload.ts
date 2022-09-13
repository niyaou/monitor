import { contextBridge } from 'electron'
import { ipcRenderer } from 'electron'

import fs from 'fs'
import { IElectronAPI } from './ui/global'

const electronAPI: IElectronAPI = {
    save: (data: string, suffix: string): void => {
        ipcRenderer.send('asynchronous-save', data, suffix);
    },
    quit: (): void => {
        ipcRenderer.send('asynchronous-close', null, null);
    },
    open: (filter?: any) => {
        const filters = []
        if (filter) {
            filters.push(filter)
        }
        const result = ipcRenderer.sendSync('asynchronous-open', filters)
        if (result !== undefined) {
            return new Promise((resolve, reject) => {
                fs.readFile(result, 'utf8', (err, data) => {
                    if (err) {
                        console.error(err)
                        return
                    }
                    let paths = result.split("\\")
                    console.log(paths[paths.length - 1])
                    paths = paths[paths.length - 1]
                    let filename = paths.split(".")[0]

                    data = JSON.parse(data)
                    let packet = { filename, data }
                    resolve(JSON.stringify(packet))
                })
            })
        } else {
            return new Promise((resolve, reject) => undefined)
        }

    },
}

ipcRenderer.on('asynchronous-save-reply', (event, arg) => {
    let filePath = arg.filePath
    if (!filePath.endsWith('.fc2')) {
        filePath += '.' + arg.suffix
    }
    fs.writeFile(filePath, arg.data, (err: any) => {
        if (err) {
            console.error(err)
            return
        }
    })
})

contextBridge.exposeInMainWorld('electronAPI', electronAPI)