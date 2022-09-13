import lodash from 'lodash'

function genArray(length: number, value: any) {
    const array = new Array(length)
    for (let index = 0; index < length; index++) {
        array[index] = lodash.cloneDeep(value)
    }
    return array;
}

const rfeChirpShapeType = {
    u32TStart: 1000,
    u32TPreSampling: 500,
    u32TPostSampling: 100,
    u32TReturn: 10000,
    u32CenterFrequency: 76.500,
    u32AcqBandwidth: 1000,
    u8TxChannelEnable: 63,
    u32aTxChannelPower: genArray(6, 10.00),
    u8aRxChannelGain: genArray(8, 7),
}



const rfeSetting = {
    u8FrontEnd: 2,
    u32RxChannelEnable: 255,
    u32SamplingFrequency: 40000,
    u16NrChirpsInFrame: 128,
    u16NrSamplesPerChirp: 128,
    u8NrChirpShapes: 1,
    straChirpShapes: genArray(8, rfeChirpShapeType),
    u8TxSwitchEnable:0x3f,
    u8aTxGain:[255,255,255,255,255,255],
}

const tef82xxAutoDriftParamType = {
    u8EnabeFreqAutoDrift: 0,
    u32FreqDriftHz: 0,
}

const tef82xxPhaseRotatorParamType = {
    u8EnablePhaseRotators: 0,
    u8aUseDDMA: genArray(6, 1),
    enuDdmaMode: 0,
    u32aDdmaInitPhase: genArray(6, 1),
    u32aDdmaPhaseUpdate: genArray(6, 2),
    enuaFinalPCGenMode: genArray(6, 0),
    u8aPhaseShiftControlSource: genArray(6, 0),
    u8EnAsyncBpskIOSampling: 0,
    u8EnAsyncQpskIOSampling: 0,
}

const tef82xxProfileOptionalParamType = {
    u32aTxPhase: genArray(6, 0),
    u8aTxBPS: genArray(6, 0),
    u8PDCBWWide: 0,
    enuaRxLPF: genArray(8, 0),
    enuaRxHPF: genArray(8, 0),
    u8VcoSel: 0,
    u8VirtualChannelNo: 0,
}

const tef82xxFrameOpParam = {
    u32SeqInterval: 0,
    u32PonDelay: 0,
    u16IsmDelay: 0,
    u8UseExtTrig: 0,
    u8ProfReset: 0,
    u8ProfModeSel: 0,
    u8aProfList: genArray(8, 0),
    u8ProfStayCnt: 1,
    u32TxPonGroupDelay: 0,
    u32RxPonGroupDelay: 0,
    u32GDelayFineControl1: 0,
    u32GDelayFineControl2: 0,
    u32GDelayFineControl3: 0,
    u32GDelayFineControl4: 0,
    u8NumSeqInBurst: 0,
    u8SafetyMontrActCtrl: 0,
    u8EnPRSafetyCheck: 0,
    u32PrSafetyStartDelay: 0,
    u8EnPRCalib: 0,
    u32PrCalibStartDelay: 0,
    u8ChirpProgressiveType: 0,
    u8SweepRstCtrl: 0,
    u8FastDischargeGSEnable: 0,
    u8FastDischargeCurrInjEnable: 0,
    u8CafcTxCalMode: 0,
    u8CafcPllLPFSel: 0,
    u8CafcPllLPFLUTSel: 1,
    u32CafcLoopBandwidth: 300000,
    straProfileOpParam: genArray(8, tef82xxProfileOptionalParamType),
    strAutoDrift: tef82xxAutoDriftParamType,
    strPhaseRotator: tef82xxPhaseRotatorParamType,
}

const defaultConfigValue = {
    u16Cmd: 0,
    u16AcqNrFrames: 1,
    strRfeSetting: rfeSetting,
    strTef82xxFrameOpParam: tef82xxFrameOpParam,
}


const getDefaultConfig = () => {
    return lodash.cloneDeep(defaultConfigValue)
}

export default getDefaultConfig