import React, { useState, useEffect, createContext, createRef } from 'react';

import { PageHeader, Button, Descriptions, Result, Space, Statistic, Modal, InputNumber, Select, Col, Row } from 'antd';
import { FolderOpenOutlined, SaveOutlined } from '@ant-design/icons';
const electronAPI = window.electronAPI
import Two from 'two.js'
import { find, findIndex } from 'lodash'
import { Switch } from 'antd';
import lodash from 'lodash'

const geometric = require("geometric");

/** 
 *  
*/
const HALF_PI = Math.PI / 2;

export interface ChirpConfig {
    u32TStart: number,
    u32TPreSampling: number,
    u32TPostSampling: number,
    u32TReturn: number,
    u32CenterFrequency: number,
    u32AcqBandwidth: number,
    u8TxChannelEnable: number,
    u32aTxChannelPower: number[],
    u8aRxChannelGain: number[],
    u32aTxPhase: number[],
    u8VcoSel: number,
}



function ProfileView(props) {
    const { index, data, setProfileValue, formName } = props

    const fieldMap = {
        "u32TStart": `strRfeSetting.straChirpShapes[${index}].u32TStart`,
        "u32TPreSampling": `strRfeSetting.straChirpShapes[${index}].u32TPreSampling`,
        "u32TPostSampling": `strRfeSetting.straChirpShapes[${index}].u32TPostSampling`,
        "u32TReturn": `strRfeSetting.straChirpShapes[${index}].u32TReturn`,
        "u32CenterFrequency": `strRfeSetting.straChirpShapes[${index}].u32CenterFrequency`,
        "u32AcqBandwidth": `strRfeSetting.straChirpShapes[${index}].u32AcqBandwidth`,
        "u8TxChannelEnable": `strRfeSetting.straChirpShapes[${index}].u8TxChannelEnable`,
        "u32aTxChannelPower": `strRfeSetting.straChirpShapes[${index}].u32aTxChannelPower`,
        "u8aRxChannelGain": `strRfeSetting.straChirpShapes[${index}].u8aRxChannelGain`,
        "u32aTxPhase": `strTef82xxFrameOpParam.straProfileOpParam[${index}].u32aTxPhase`,
        "u8VcoSel": `strTef82xxFrameOpParam.straProfileOpParam[${index}].u8VcoSel`,
    }

    const formValue: ChirpConfig = {
        u32TStart: 0,
        u32TPreSampling: 0,
        u32TPostSampling: 0,
        u32TReturn: 0,
        u32CenterFrequency: 0,
        u32AcqBandwidth: 0,
        u8TxChannelEnable: 0,
        u32aTxChannelPower: [0, 0, 0, 0, 0, 0],
        u8aRxChannelGain: [0, 0, 0, 0, 0, 0, 0, 0],
        u8VcoSel: 0,
        u32aTxPhase: [0, 0, 0, 0, 0, 0],
    }
    
    Object.keys(fieldMap).forEach(key => {
        formValue[`${key}`] = lodash.get(data, fieldMap[key])
    });
    // console.log(defaultValue)

    const { Option } = Select;
    const vcoSel = [{ key: '1GH_VCO', value: 0 }, { key: '2GHZ_VCO', value: 1 }, { key: '4GHz_VCO', value: 2 },]

    const aRxChannelGain = [{ key: '27db', value: 1 }, { key: '30db', value: 2 }, { key: '33db', value: 3 }, { key: '36db', value: 4 }, { key: '39db', value: 5 }, { key: '42db', value: 6 }, { key: '45db', value: 7 },]
    const [visible, setVisible] = useState(false);
    const [vcoValue, setVcoValue] = useState(formValue.u8VcoSel);

    const [channelEnable, setChannelEnable] = useState(formValue.u8TxChannelEnable);
    const [channelPower, setChannelPower] = useState(formValue.u32aTxChannelPower);
    const [channelGain, setChannelGain] = useState(formValue.u8aRxChannelGain);
    const [phase, setPhase] = useState(formValue.u32aTxPhase);

    const [value, setValue] = useState(0);
    const [displayText, setDisplayText] = useState('');
    const [unitText, setUnitText] = useState('');
    // const [config, setConfig] = useState<ChirpConfig>(_config);
    const [config, setConfig] = useState<ChirpConfig>(formValue);
    const [callbackName, setCallbackName] = useState(null);
    const [switchs, setSwitchs] = useState([]);
    const [two, setTwo] = useState<Two>();

    console.log(data)
    useEffect(() => {
        let _conf = { ...config }
        _conf.u8TxChannelEnable = channelEnable
        _conf.u32aTxChannelPower = channelPower
        _conf.u8aRxChannelGain = channelGain
        _conf.u32aTxPhase = phase
        _conf.u8VcoSel = vcoValue
        setConfig(_conf)
    }, [channelEnable, channelPower, channelGain, channelGain, phase, vcoValue])


    const handleChange = (value: number) => {
        console.log(`selected ${value}`);
        setVcoValue(value)
    };
    const handleOk = (e: React.MouseEvent<HTMLElement>) => {
        setVisible(false);
        let _conf = { ...config }
        _conf[callbackName] = value
        // console.log("🚀 ~----------- config", _conf)
        setConfig(_conf)
    };

    const handleCancel = (e: React.MouseEvent<HTMLElement>) => {
        setVisible(false);
    };


    const makeIndicatorLine = (two, start, end, dashed = false, startArr = false, endArr = false, deltaS = { x: 0, y: 0 }, deltaE = { x: 0, y: 0 }) => {
        let group = two.makeGroup();
        let line = two.makeLine(start.x, start.y, end.x, end.y);
        line.fill = 'rgb(234, 60, 50)';
        line.linewidth = 1;
        let deltaX, deltaY = 0
        deltaX = end.x - start.x
        deltaY = end.y - start.y
        if (dashed) {
            line.dashes = [5, 1]
        }

        group.add(line)
        if (startArr) {
            let anther = makeTriangle(two)
            // anther.rotation = geometric.lineAngle([[start.x,start.y], [end.x, end.y]]) ; 
            anther.translation.set(start.x - deltaS.x, start.y - deltaS.y)
            anther.rotation = Math.atan2(-deltaY, -deltaX) - Math.PI / 2;
            group.add(anther)
        }
        if (endArr) {

            let triangle = makeTriangle(two)
            triangle.translation.set(end.x - deltaE.x, end.y - deltaE.y)
            triangle.rotation = Math.atan2(-deltaY, -deltaX) + Math.PI / 2;
            group.add(triangle)
        }
        return group
    }

    const makeSpecialLine = (two, start, end, dashed = false, startArr = false, endArr = false) => {
        let group = two.makeGroup();
        let line = two.makeLine(start.x, start.y, end.x, end.y);
        line.stroke = 'rgb(	218,165,32)';

        line.linewidth = 1;
        let deltaX, deltaY = 0
        deltaX = end.x - start.x
        deltaY = end.y - start.y
        if (dashed) {
            line.dashes = [5, 1]
        }
        let triangle = two.makeCircle(end.x, end.y, 4)
        triangle.stroke = "rgb(218,165,32)"
        triangle.fill = "rgb(218,165,32)"
        group.add(triangle)
        return group
    }


    const makeTriangle = (two) => {
        let size = 8
        var tri = two.makePath(- size / 2, 0, size / 2, 0, 0, size);
        tri.fill = 'rgb(0, 0, 0)';
        return tri;
    }


    const MARGIN_Y = -120
    const DWEELS_X = 110
    const DWEELE_X = 300
    const DWEELE_Y = 500 + MARGIN_Y
    const STOP_X = 800
    const STOP_Y = 180 + MARGIN_Y
    const FOLLOW_Y = 240 + MARGIN_Y
    const CHRISPE_X = 960
    const EFFECT_Y = 320 + MARGIN_Y
    const MIDDLE_Y = 420 + MARGIN_Y
    const AXIOSB_Y = 600 + MARGIN_Y
    const AXIOST_Y = 130 + MARGIN_Y
    const verticalHead_y = 630 + MARGIN_Y
    const verticalTail_y = 650 + MARGIN_Y
    const BUTTOMV_Y = 700 + MARGIN_Y
    const BUTTOM_Y = 730 + MARGIN_Y
    const drawPoint = {
        axisT: { x: 50, y: AXIOST_Y }, axiosB: { x: 50, y: AXIOSB_Y }, axisR: { x: 1050, y: AXIOSB_Y },
        interptStop: { x: 50, y: STOP_Y }, interptF: { x: 50, y: FOLLOW_Y },
        interptStopE: { x: STOP_X, y: STOP_Y },
        interptJ: { x: 50, y: MIDDLE_Y }, interptStart: { x: 50, y: DWEELE_Y },
        interptFE: { x: 50, y: FOLLOW_Y }, interptE: { x: 50, y: EFFECT_Y }, interptEE: { x: 0, y: EFFECT_Y },
        interptJE: { x: 0, y: MIDDLE_Y }, interptJEE: { x: 0, y: MIDDLE_Y },
        dwellS: { x: DWEELS_X, y: DWEELE_Y }, dwellE: { x: DWEELE_X, y: DWEELE_Y }, resetS: { x: STOP_X, y: DWEELE_Y },
        chrispE: { x: CHRISPE_X, y: DWEELE_Y }, dwellVS: { x: DWEELS_X, y: verticalHead_y },
        dwellVE: { x: DWEELE_X, y: verticalHead_y }, settleE: { x: 0, y: verticalHead_y }, jumpE: { x: 0, y: verticalHead_y },
        interptStopEV: { x: STOP_X, y: verticalHead_y }, resetVE: { x: CHRISPE_X, y: verticalHead_y },
        dwellVV: { x: DWEELS_X, y: BUTTOM_Y }, chrispVS: { x: DWEELS_X, y: BUTTOMV_Y }, chrispVE: { x: CHRISPE_X, y: BUTTOMV_Y }, chrispVV: { x: CHRISPE_X, y: BUTTOM_Y },
        dwellVVE: { x: DWEELE_X, y: verticalTail_y }, settleVE: { x: 0, y: verticalTail_y },
        jumpVE: { x: 0, y: verticalTail_y }, interptStopEVV: { x: STOP_X, y: verticalTail_y },
    }

    const getInterptY = (startPoint, endPoint, percent) => {
        let interptEE = geometric.lineInterpolate([[startPoint.x, startPoint.y], [endPoint.x, endPoint.y]])
        interptEE = interptEE(percent)
        return interptEE
    }



    const typeText = (two, position, content, deltaX, deltaY, needSwitch = false, callbackName = '', displayText = '', unitText = '') => {
        let group = two.makeGroup();
        let line = two.makeText(content);
        line.translation.set(position.x + deltaX, position.y + deltaY)
        let lastSwitch = switchs
        group.add(line)
        if (needSwitch) {
            // console.log("🚀 ~ file: profile.tsx ~ line 170 ~ typeText ~ config", config)
            let _switch = {
                _conf: config, callbackName, position: line.position, callback: () => {
                    // console.log("🚀 ~ file: profile.tsx ~ line 170 ~ typeText ~ config", config)
                    setValue(config[callbackName])
                    setDisplayText(displayText)
                    setUnitText(unitText)
                    setVisible(true)
                    setCallbackName(callbackName)
                }
            }

            let hit = findIndex(lastSwitch, ['callbackName', callbackName])
            if (hit > -1) {
                lastSwitch[hit] = _switch
            } else {
                lastSwitch.push(_switch)
            }

            setSwitchs(lastSwitch)
            // console.log("🚀 ~ file: profile.tsx ~ line 176 ~ switchs.push ~ switchs", switchs)
        }
        return line
    }

    let _click, _move = null

    useEffect(() => {
        setConfig(formValue)

        let artboard = document.querySelector('#artboard' + formName);
        let _two = new Two({
            type: Two.Types.svg,
            height: 710,
            width: 1080
        }).appendTo(artboard as HTMLBaseElement);
        setTwo(_two)
        _move = (e) => {
            let matched = false;
            let position = new Two.Vector(e.offsetX, e.offsetY)
            for (let i = 0; i < switchs.length; i++) {
                const dist = Two.Vector.distanceBetween(switchs[i].position, position);
                //   console.log("🚀 ~ file: profile.tsx ~ line 232 ~ useEffect ~ freq",e.offsetX,e.offsetY,dist, freq.position)
                if (dist < 32) {
                    matched = true
                    break
                }
            }
            if (matched) {
                _two.renderer.domElement.style.cursor = 'pointer';
            } else {
                _two.renderer.domElement.style.cursor = 'auto';
            }
        }

        _click = (e) => {
            let matched = false;
            let position = new Two.Vector(e.offsetX, e.offsetY)
            for (let i = 0; i < switchs.length; i++) {
                const dist = Two.Vector.distanceBetween(switchs[i].position, position);
                //   console.log("🚀 ~ file: profile.tsx ~ line 232 ~ useEffect ~ freq",e.offsetX,e.offsetY,dist, freq.position)
                if (dist < 32) {
                    matched = true
                    switchs[i].callback()
                    break

                }
            }
        }
        window.addEventListener('mousemove', _move)
        window.addEventListener('mousedown', _click)
        // return () => {
        //     console.log("11111111111111111111111")
        // };
    }, [])


    const getReformValue = () => {
        const reformValue = {}
        Object.keys(config).forEach(key => {
            lodash.set(reformValue, fieldMap[key], config[key])
        });
        return reformValue
    }

    useEffect(() => {
        if (typeof setProfileValue === 'function') {
            // callback(config)
            // console.log("callback")
            setProfileValue(formName, getReformValue());
        }
        if (!two) {
            return
        }
        setSwitchs([])
        two.clear()

        // arrows.translation.set(two.width / 2, two.height / 2);
        // let x = 210
        // let y = 300
        makeIndicatorLine(two, drawPoint.axisT, drawPoint.axiosB, false, true, false)
        makeIndicatorLine(two, drawPoint.axiosB, drawPoint.axisR, false, false, true)
        makeIndicatorLine(two, drawPoint.interptStop, drawPoint.interptStopE, true, false, false)
        makeIndicatorLine(two, drawPoint.dwellE, drawPoint.interptStopE, false, false, false)
        let interptFE = geometric.lineInterpolate([[drawPoint.dwellE.x, drawPoint.dwellE.y], [drawPoint.interptStopE.x, drawPoint.interptStopE.y]])
        interptFE = interptFE((drawPoint.dwellE.y - drawPoint.interptF.y) / (drawPoint.dwellE.y - drawPoint.interptStopE.y))
        drawPoint.interptFE.x = interptFE[0]
        drawPoint.interptFE.y = interptFE[1]
        drawPoint.interptJEE.x = interptFE[0]
        drawPoint.jumpVE.x = interptFE[0]
        drawPoint.jumpE.x = interptFE[0]
        makeIndicatorLine(two, drawPoint.interptF, drawPoint.interptFE, true, false, false)
        let interptEE = geometric.lineInterpolate([[drawPoint.dwellE.x, drawPoint.dwellE.y], [drawPoint.interptStopE.x, drawPoint.interptStopE.y]])
        interptEE = interptEE((drawPoint.dwellE.y - drawPoint.interptE.y) / (drawPoint.dwellE.y - drawPoint.interptStopE.y))
        drawPoint.interptEE.x = interptEE[0]
        drawPoint.interptEE.y = interptEE[1]
        makeSpecialLine(two, drawPoint.interptE, drawPoint.interptEE, true, false, false)
        makeIndicatorLine(two, drawPoint.interptJ, drawPoint.interptJEE, true, false, false)
        makeIndicatorLine(two, drawPoint.interptStart, drawPoint.chrispE, true, false, false)
        makeIndicatorLine(two, drawPoint.interptStopE, drawPoint.chrispE, false, false, false)
        makeIndicatorLine(two, drawPoint.interptStopE, drawPoint.resetS, false, true, true, { x: 0, y: -5 }, { x: 0, y: 5 })
        makeIndicatorLine(two, drawPoint.interptFE, drawPoint.interptJEE, false, true, true, { x: 0, y: -5 }, { x: 0, y: 5 })
        makeIndicatorLine(two, drawPoint.dwellS, drawPoint.dwellE, false, false, false)
        makeIndicatorLine(two, drawPoint.dwellS, drawPoint.dwellVV, true, false, false)
        makeIndicatorLine(two, drawPoint.chrispE, drawPoint.chrispVV, true, false, false)
        makeIndicatorLine(two, drawPoint.resetS, drawPoint.interptStopEVV, true, false, false)
        makeIndicatorLine(two, drawPoint.interptJEE, drawPoint.jumpVE, true, false, false)
        let p = getInterptY(drawPoint.dwellE, drawPoint.interptStopE, (drawPoint.dwellE.y - drawPoint.interptJ.y) / (drawPoint.dwellE.y - drawPoint.interptStopE.y))
        drawPoint.interptJE.x = p[0]
        drawPoint.settleVE.x = p[0]
        drawPoint.settleE.x = p[0]
        makeIndicatorLine(two, drawPoint.interptJE, drawPoint.settleVE, true, false, false)
        makeIndicatorLine(two, drawPoint.dwellE, drawPoint.dwellVVE, true, false, false)
        makeIndicatorLine(two, drawPoint.dwellVS, drawPoint.dwellVE, false, true, true, { x: -5, y: 0 }, { x: 5, y: 0 })
        makeIndicatorLine(two, drawPoint.dwellVE, drawPoint.settleE, false, true, true, { x: -5, y: 0 }, { x: 5, y: 0 })
        makeIndicatorLine(two, drawPoint.settleE, drawPoint.jumpE, false, true, true, { x: -5, y: 0 }, { x: 5, y: 0 })
        makeIndicatorLine(two, drawPoint.jumpE, drawPoint.interptStopEV, false, true, true, { x: -5, y: 0 }, { x: 5, y: 0 })
        makeIndicatorLine(two, drawPoint.interptStopEV, drawPoint.resetVE, false, true, true, { x: -5, y: 0 }, { x: 5, y: 0 })
        makeIndicatorLine(two, drawPoint.chrispVS, drawPoint.chrispVE, false, true, true, { x: -5, y: 0 }, { x: 5, y: 0 })

        typeText(two, drawPoint.interptStop, "Stop Freq.    -", 90, -10)
        typeText(two, drawPoint.interptE, "Effective Center Freq.  ", 120, 30)
        typeText(two, drawPoint.interptE, `${config.u32CenterFrequency}MHz`, 240, 30, true, 'u32CenterFrequency', 'Effective Center Freq', 'MHz')
        typeText(two, drawPoint.interptStart, "Start Freq.     MHz", 90, -20)
        // typeText(two, drawPoint.axiosB, "Loop BandWidth", 2, 10)
        typeText(two, drawPoint.dwellVS, `Dwell Time ${config.u32TStart}ns`, 80, -10, true, 'u32TStart', 'Dwell Time', 'ns')
        typeText(two, drawPoint.dwellVE, `Settle Time ${config.u32TPreSampling}ns`, 60, -10, true, 'u32TPreSampling', 'Settle Time', 'ns')
        typeText(two, drawPoint.settleE, "Sample Time -ns", 90, -10,)
        typeText(two, drawPoint.jumpE, "Jump back Time", 50, 30)
        typeText(two, drawPoint.jumpE, `${config.u32TPostSampling}ns`, 35, -12, true, 'u32TPostSampling', 'Jump back Time', 'ns')
        typeText(two, drawPoint.interptStopEV, `Reset Time ${config.u32TReturn}ns`, 90, 20, true, 'u32TReturn', 'Reset Time', 'ns')
        typeText(two, drawPoint.settleVE, "chirp Time -us", 90, 30)
        typeText(two, drawPoint.interptEE, "Effective BandWidth ", 90, 20)
        typeText(two, drawPoint.interptEE, `${config.u32AcqBandwidth}MHZ`, 190, 20, true, 'u32AcqBandwidth', 'Effective BandWidth', 'MHz')
        let title = typeText(two, drawPoint.axiosB, "frequency ", -20, -300)
        title.size = 24
        title.width = 800
        title.rotation = -HALF_PI
        two.render();
    }, [config, two]);

    const saveConfig = () => {
        electronAPI.save(JSON.stringify(config), 'fc2profile')
    };

    const loadConfig = () => {
        electronAPI.open({ name: 'Fc2', extensions: ['fc2profile'] }).then(e => {
            // console.log(e)
            const _conf = JSON.parse(e)
            setChannelEnable(_conf.u8TxChannelEnable)
            setChannelPower(_conf.u32aTxChannelPower)
            setChannelGain(_conf.u8aRxChannelGain)
            setPhase(_conf.u32aTxPhase)
            setVcoValue(_conf.u8VcoSel)
            setConfig(_conf)

            setProfileValue(formName, getReformValue());
        })
    }

    useEffect(() => {
        const _conf = formValue
        setChannelEnable(_conf.u8TxChannelEnable)
        setChannelPower(_conf.u32aTxChannelPower)
        setChannelGain(_conf.u8aRxChannelGain)
        setPhase(_conf.u32aTxPhase)
        setVcoValue(_conf.u8VcoSel)
        // setConfig(_conf)

        // setProfileValue(formName, getReformValue());
    }, [data])

    return (<>
        <PageHeader
            className="site-page-header"
            title={formName}
            extra={
                [
                    <Button
                        key="extra-load6"
                        icon={<FolderOpenOutlined />}
                        onClick={() => loadConfig()}
                    >
                        加载配置
                    </Button>,
                    <Button
                        key="extra-load5"
                        icon={<SaveOutlined />}
                        onClick={() => saveConfig()}
                    >
                        保存配置
                    </Button>,
                ]
            }
        />
        <Row gutter={2}>
            <Col span={17}>
                <div id={"artboard" + formName} style={{ overflowX:'auto',maxHeight:"770px", borderWidth: 1, borderStyle: "solid", borderColor: "#c2c2c2" }} >
                    <Row justify="start" align="middle" style={{ marginTop: '20px' }}>
                        <Col span={4} offset={1}><div>Chirp波形参数设置</div></Col>
                        <Col span={1} offset={11}><div>VcoSel</div></Col>
                        <Col span={3} offset={0}>
                            <Select
                                value={vcoValue}
                                dropdownStyle={{ background: "rbg(240,242,245)" }}
                                style={{ width: 120, background: '#00000000' }}
                                bordered={false}
                                onChange={handleChange}>
                                {vcoSel.map((v, i) => { return (<Option key={i} value={v.value} style={{ background: '#00000000' }}>{v.key}</Option>) })}
                            </Select>
                        </Col>
                    </Row>
                </div>
            </Col>
            <Col span={7}>
                <Row justify="center" align="middle" style={{ borderWidth: 1, borderStyle: "solid", borderColor: "#c2c2c2", paddingTop: 5,height:'380px' }} >
                    <Col span={5} offset={1}><div>发射天线设置</div></Col>
                    <Col span={6}><div>功率</div></Col>
                    <Col span={6}><div>相位</div></Col>  
                    <Col span={6}><div>发射开关</div></Col>
                    <Row justify="start" align="middle" style={{ margin: 3 }} >
                        {[1, 2, 3, 4, 5, 6].map((v, i) => {

                            // console.log("🚀 ~ file: profile.tsx ~ line 386 ~ {[1,2,3,4,5,6].map ~ {channelEnable&&(0x1<<i)===1", channelEnable,
                            //   0x1<<i,i, channelEnable&(1<<i))
                            return (<Col span={24} key={`receive-${i}`} >
                                <Row justify="start" align="middle" style={{ margin: 15 }} >
                                    <Col span={6}>  <img src={require('../../../../asset/signal-c.png')} style={{ width: 24, height: 24 }} />
                                        天线_{v}</Col>
                                    <Col span={6}> <InputNumber size="small" min={0} max={100000} value={channelPower[i]}
                                        onChange={(v) => {
                                            let power = [].concat(channelPower)
                                            power[i] = v
                                            setChannelPower(power)
                                            // console.log("🚀 ~ file: profile.tsx ~ line 399 ~ .map ~ v", v)
                                        }}
                                        bordered={false} /></Col>
                                    <Col span={6}> <InputNumber size="small" min={0} max={100000}
                                        value={phase[i]}
                                        onChange={(v) => {
                                            let _phase = [].concat(phase)
                                            _phase[i] = v
                                            setPhase(_phase)
                                            //  console.log("🚀 ~ file: profile.tsx ~ line 399 ~ .map ~ v", v)
                                        }}

                                        bordered={false} /></Col>
                                    <Col span={6}> <Switch checkedChildren="开启" unCheckedChildren="关闭"
                                        onChange={(checked, event) => {
                                            // console.log("🚀 ~ file: profile.tsx ~ line 369 ~ ProfileView ~ checked,event",i, checked,event)

                                            config.u8TxChannelEnable = channelEnable
                                            let enable = checked ? channelEnable | (1 << i) : channelEnable & (~(1 << i))
                                            setChannelEnable(enable)

                                        }}
                                        checked={(channelEnable & (1 << i)) === (1 << i)} /></Col>
                                </Row>
                            </Col>)
                        })}
                    </Row>
                </Row>
                <Row justify="center" align="middle" style={{ borderWidth: 1, borderStyle: "solid", borderColor: "#c2c2c2", marginTop: 8 ,height:'380px' }} >
                    <Col span={6} offset={1}>接收天线设置</Col>
                    <Col span={6}>功率</Col>
                    <Col span={5}></Col>
                    <Col span={6}>功率</Col>
                    <Row justify="start" align="middle" style={{ margin: 5 }} >
                        {[1, 2, 3, 4, 5, 6, 7, 8].map((v, i) => {
                            return (<Col span={12} key={`receive-${i}`} >
                                <Row justify="start" align="middle" style={{ margin: 15 }} >
                                    <Col span={12}>  <img src={require('../../../../asset/signal-r.png')} style={{ width: 20, height: 20 }} />
                                        天线_{v}</Col>
                                    <Col span={12}>
                                        <Select
                                            value={channelGain[i]}
                                            dropdownStyle={{ background: "rbg(240,242,245)" }}
                                            style={{ width: 120, background: '#00000000' }}
                                            bordered={false}
                                            onChange={
                                                (v, e) => {
                                                    // console.log("🚀 ~ file: profile.tsx ~ line 454 ~ {[1,2,3,4,5,6,7,8].map ~ v,e", v,e)
                                                    let _gain = [].concat(channelGain)
                                                    _gain[i] = v
                                                    setChannelGain(_gain)
                                                }
                                            }>
                                            {aRxChannelGain.map((v, i) => { return (<Option key={i} value={v.value} style={{ background: '#00000000' }}>{v.key}</Option>) })}
                                        </Select>
                                    </Col>

                                </Row>
                            </Col>)
                        })}
                    </Row>
                    <Row justify="center" align="middle" style={{ marginTop: '20px' }}>

                    </Row>

                </Row>
            </Col>
        </Row>

        <Modal
            centered={true}
            visible={visible}
            onOk={handleOk}
            onCancel={handleCancel}
        >
            <Space>
                {/* <InputNumber size="large" min={1} max={100000} defaultValue={3} /> */}
                <div>{displayText}</div>
                {/* <InputNumber min={1} max={100000} defaultValue={3}  /> */}
                <InputNumber size="small"
                    defaultValue={3}
                    value={value}
                    onChange={(v) => {
                        setValue(v)
                    }}
                />
                <div> {unitText}</div>

            </Space>
        </Modal>
    </>);
}

export default ProfileView;