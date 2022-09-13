export interface IElectronAPI {
    save: (data: string, suffix: string) => void,
    quit: () => void,
    open: (filter?: any) => Promise<string>,
}

declare global {
    interface Window {
        electronAPI: IElectronAPI
    }
}