import React, { useState, useEffect } from 'react';
import { Row, Col, Tabs, Form, Input, Switch, Radio, Checkbox, Select, InputNumber, FormInstance, PageHeader, Button, Descriptions, Result, Space, Statistic, Modal, } from 'antd';
import { IForm } from './types'
const electronAPI = window.electronAPI
import lodash from 'lodash'
import { ChirpConfig } from './profile'
const { Option } = Select;
import { cloneDeep, find, findIndex } from 'lodash'

function getBit(data: number, num: number) {
  const byteData = data;
  const byteArr = new Array(8);
  byteArr[0] = (byteData & 0x01) == 0x01 ? 1 : 0;
  byteArr[1] = (byteData & 0x02) == 0x02 ? 1 : 0;
  byteArr[2] = (byteData & 0x04) == 0x04 ? 1 : 0;
  byteArr[3] = (byteData & 0x08) == 0x08 ? 1 : 0;
  byteArr[4] = (byteData & 0x10) == 0x10 ? 1 : 0;
  byteArr[5] = (byteData & 0x20) == 0x20 ? 1 : 0;
  byteArr[6] = (byteData & 0x40) == 0x40 ? 1 : 0;
  byteArr[7] = (byteData & 0x80) == 0x80 ? 1 : 0;
  return byteArr[num];
}

function PointSet(props) {
  const { callback, config: _config } = props
  const { index, data, setProfileValue, formName } = props

  const [form] = Form.useForm();

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
    "strTef82xxFrameOpParam.u32CafcLoopBandwidth": `strTef82xxFrameOpParam.u32CafcLoopBandwidth`,
    "strTef82xxFrameOpParam.u8CafcPllLPFLUTSel": `strTef82xxFrameOpParam.u8CafcPllLPFLUTSel`,
  }

  const formValue: any = {}

  Object.keys(fieldMap).forEach(key => {
    formValue[`${key}`] = lodash.get(data, fieldMap[key])
  });

  const { Option } = Select;
  const [vcoValue, setVcoValue] = useState(formValue.u8VcoSel);
  const [channelEnable, setChannelEnable] = useState(formValue.u8TxChannelEnable);
  const [u8TxSwitchEnable, setU8TxSwitchEnable] = useState(lodash.get(data, "strRfeSetting.u8TxSwitchEnable"));
  const [u8aTxGain, setU8aTxGain] = useState(lodash.get(data, "strRfeSetting.u8aTxGain"));
  const [u32CenterFrequency, setU32CenterFrequency] = useState(formValue.u32CenterFrequency);
  const [channelPower, setChannelPower] = useState(formValue.u32aTxChannelPower);
  const [channelGain, setChannelGain] = useState(formValue.u8aRxChannelGain);
  const [phase, setPhase] = useState(formValue.u32aTxPhase);

  const [config, setConfig] = useState<any>(formValue);
  const [u8CafcPllLPFLUTSel, setU8CafcPllLPFLUTSel] = useState(lodash.get(data, "strTef82xxFrameOpParam.u8CafcPllLPFLUTSel"));
  const [u32CafcLoopBandwidth, setU32CafcLoopBandwidth] = useState(lodash.get(data, "strTef82xxFrameOpParam.u32CafcLoopBandwidth"));

  useEffect(() => {
    let _conf = { ...formValue }
    _conf.u8TxChannelEnable = channelEnable
    _conf.u32aTxChannelPower = channelPower
    _conf.u8aRxChannelGain = channelGain
    _conf.u32aTxPhase = phase
    _conf.u8VcoSel = vcoValue
    _conf.u32CenterFrequency = u32CenterFrequency

    setConfig(_conf)

    let _chirp = data.strRfeSetting.straChirpShapes

    _chirp[index] = _conf
    setProfileValue('', { strRfeSetting: { straChirpShapes: _chirp } })
  }, [channelEnable, channelPower, channelGain, channelGain, phase, vcoValue, u8TxSwitchEnable, u32CenterFrequency])


  useEffect(() => {
    setProfileValue('strTef82xxFrameOpParam.u8CafcPllLPFLUTSel', { strTef82xxFrameOpParam: { u8CafcPllLPFLUTSel } })
  }, [u8CafcPllLPFLUTSel])

  useEffect(() => {
    setProfileValue('strTef82xxFrameOpParam.u32CafcLoopBandwidth', { strTef82xxFrameOpParam: { u32CafcLoopBandwidth } })
  }, [u32CafcLoopBandwidth])

  useEffect(() => {
    setProfileValue('strRfeSetting.u8TxSwitchEnable', { strRfeSetting: { u8TxSwitchEnable } })
  }, [u8TxSwitchEnable])

  useEffect(() => {
    setProfileValue('strRfeSetting.u8aTxGain', { strRfeSetting: { u8aTxGain } })
  }, [u8aTxGain])



  useEffect(() => {
    // console.log("🚀 ~ file: pointset.tsx ~ line 33 ~ PointSet ~ data", cloneDeep(data), index)
    let _conf = data.strRfeSetting.straChirpShapes[0]
    setChannelEnable(_conf.u8TxChannelEnable)
    setChannelPower(_conf.u32aTxChannelPower)
    setChannelGain(_conf.u8aRxChannelGain)
    setPhase(_conf.u32aTxPhase)
    setVcoValue(_conf.u8VcoSel)

    let U32CenterFrequency = data.strRfeSetting.straChirpShapes[0].u32CenterFrequency
    let U8CafcPllLPFLUTSel = lodash.get(data, "strTef82xxFrameOpParam.u8CafcPllLPFLUTSel")
    let U32CafcLoopBandwidth = lodash.get(data, "strTef82xxFrameOpParam.u32CafcLoopBandwidth")
    setU32CenterFrequency(U32CenterFrequency);
    setU8CafcPllLPFLUTSel(U8CafcPllLPFLUTSel)
    setU32CafcLoopBandwidth(U32CafcLoopBandwidth)
    setU8TxSwitchEnable(lodash.get(data, "strRfeSetting.u8TxSwitchEnable"))

    form.resetFields()
  }, [data])


  const getReformValue = () => {
    const reformValue = {}
    Object.keys(config).forEach(key => {
      lodash.set(reformValue, fieldMap[key], config[key])
    });
    return reformValue
  }



  // const saveConfig = () => {
  //   electronAPI.save(JSON.stringify(config), 'fc2profile')
  // };


  // const loadConfig = () => {
  //   electronAPI.open({ name: 'Fc2', extensions: ['fc2profile'] }).then(e => {
  //     // console.log(e)
  //     const _conf = JSON.parse(e)
  //     setChannelEnable(_conf.u8TxChannelEnable)
  //     setChannelPower(_conf.u32aTxChannelPower)
  //     setChannelGain(_conf.u8aRxChannelGain)
  //     setPhase(_conf.u32aTxPhase)
  //     setVcoValue(_conf.u8VcoSel)
  //     setConfig(_conf)
  //     setProfileValue(formName, getReformValue());
  //   })
  // }

  useEffect(() => {

    const _conf = formValue
    setChannelEnable(_conf.u8TxChannelEnable)
    setChannelPower(_conf.u32aTxChannelPower)
    setChannelGain(_conf.u8aRxChannelGain)
    setPhase(_conf.u32aTxPhase)
    setVcoValue(_conf.u8VcoSel)
    setConfig(_conf)
    // setProfileValue(formName, getReformValue());
  }, [])

  return (<>

    <Row gutter={2}>
      <Col span={8}>
        <Form
          form={form}
          name={formName}
          labelCol={{ span: 9 }}
          wrapperCol={{ span: 15 }}
          initialValues={formValue}

          autoComplete="off"
        >
          <Form.Item

            label="LUT表"
            name="strTef82xxFrameOpParam.u8CafcPllLPFLUTSel"
            tooltip=" 范围:[1,65535]"
            rules={[{
              required: true,
              // type: 'number', min: 0, max: 65535 
            }]}
          >
            <Select
              dropdownStyle={{ background: "rbg(240,242,245)" }}
              style={{ width: '100%', background: '#00000000' }}
              bordered={false}
              value={u8CafcPllLPFLUTSel}
              onChange={(v) => {
                console.log("🚀 ~ file: pointset.tsx ~ line 214 ~ PointSet ~ v", v)
                setU8CafcPllLPFLUTSel(v)
              }}
            >
              {['1G_LUT',
                '5G_NARROW_LUT',
                '5G_WIDE_LUT'
              ].map((v, i) => { return (<Option key={`select-lut-${i}`} value={i} style={{ background: '#00000000' }}>{v}</Option>) })}
            </Select>
            {/* <InputNumber style={{ width: '100%' }} min={1} max={65535} /> */}
          </Form.Item>
          <Form.Item
            label="PLL环路带宽"
            name="strTef82xxFrameOpParam.u32CafcLoopBandwidth"
            tooltip={`范围:[${u8CafcPllLPFLUTSel === 0 ? 200000 : u8CafcPllLPFLUTSel === 1 ? 250000 : 300000} ~ ${u8CafcPllLPFLUTSel === 0 ? 1600000 : u8CafcPllLPFLUTSel === 1 ? 1650000 : 1700000}]`}
            rules={[{
              required: true,
              // type: 'number', min: 0, max: 65535 
            }]}
          >


            <InputNumber
              value={u32CafcLoopBandwidth}
              onChange={(v) => {
                console.log("🚀 ~ file: pointset.tsx ~ line 247 ~ ].map ~ v", v)
                setU32CafcLoopBandwidth(v)
              }}
              style={{ width: '100%' }} min={u8CafcPllLPFLUTSel === 0 ? 200000 : u8CafcPllLPFLUTSel === 1 ? 250000 : 300000}
              max={u8CafcPllLPFLUTSel === 0 ? 1600000 : u8CafcPllLPFLUTSel === 1 ? 1650000 : 1700000} />
          </Form.Item>

          <Form.Item
            label="Effective Center Freq"
            name="u32CenterFrequency"
            tooltip={`范围:[76.000 ~ 81.000]`}
            rules={[{
              required: true,
              // type: 'number', min: 0, max: 65535 
            }]}
          >


            <InputNumber style={{ width: '100%' }}
              value={u32CenterFrequency}
              onChange={(v) => {
                console.log("🚀 ~ file: pointset.tsx ~ line 267 ~ ].map ~ v", v)
                setU32CenterFrequency(v)
              }}
              min={76}
              max={81} />
          </Form.Item>
        </Form>
      </Col>
      <Col span={12} offset={4}>
        <Row justify="start" align="middle" style={{ borderWidth: 1, borderStyle: "solid", borderColor: "#c2c2c2", paddingTop: 5, height: '380px' }} >
          <Col span={4} offset={1}><div>发射天线设置</div></Col>
          <Col span={3}><div>功率</div></Col>
          <Col span={3}><div>相位</div></Col>
          <Col span={3}><div>天线增益</div></Col>
          <Col span={4}><div>发射开关</div></Col>
          <Col span={4}><div>直流电源</div></Col>
          <Row justify="start" align="middle" style={{ margin: 3 }} >
            {[1, 2, 3, 4, 5, 6].map((v, i) => {

              // console.log("🚀 ~ file: profile.tsx ~ line 386 ~ {[1,2,3,4,5,6].map ~ {channelEnable&&(0x1<<i)===1", channelEnable,
              //   0x1<<i,i, channelEnable&(1<<i))
              return (<Col span={24} key={`receive-${i}`} >
                <Row justify="start" align="middle" style={{ margin: 15 }} >
                  <Col span={5}>  <img src={require('../../../../asset/signal-c.png')} style={{ width: 24, height: 24 }} />
                    天线_{v}</Col>
                  <Col span={3}> <InputNumber size="small" min={0} max={100000} value={channelPower[i]}
                    onChange={(v) => {
                      let power = [].concat(channelPower)
                      power[i] = v
                      setChannelPower(power)
                      // console.log("🚀 ~ file: profile.tsx ~ line 399 ~ .map ~ v", v)
                    }}
                    bordered={false} /></Col>
                  <Col span={3}> <InputNumber size="small" min={0} max={100000}
                    value={phase[i]}
                    onChange={(v) => {
                      let _phase = [].concat(phase)
                      _phase[i] = v
                      setPhase(_phase)
                      //  console.log("🚀 ~ file: profile.tsx ~ line 399 ~ .map ~ v", v)
                    }}

                    bordered={false} /></Col>
                  <Col span={3}> <InputNumber size="small" min={0} max={255}
                    value={u8aTxGain[i]}
                    onChange={(v) => {
                      let _u8aTxGain = [].concat(u8aTxGain)
                      _u8aTxGain[i] = v
                      setU8aTxGain(_u8aTxGain)
                      console.log("🚀 ~ file: pointset.tsx ~ line 292 ~ {[1,2,3,4,5,6].map ~ u8aTxGain", _u8aTxGain)
                      //  console.log("🚀 ~ file: profile.tsx ~ line 399 ~ .map ~ v", v)
                    }}

                    bordered={false} /></Col>
                  <Col span={4}> <Switch checkedChildren="开启" unCheckedChildren="关闭"
                    onChange={(checked, event) => {
                      // console.log("🚀 ~ file: profile.tsx ~ line 369 ~ ProfileView ~ checked,event",i, checked,event)

                      config.u8TxChannelEnable = channelEnable
                      let enable = checked ? channelEnable | (1 << i) : channelEnable & (~(1 << i))
                      setChannelEnable(enable)

                    }}
                    checked={(channelEnable & (1 << i)) === (1 << i)} /></Col>
                  <Col span={4}> <Switch checkedChildren="开启" unCheckedChildren="关闭"
                    onChange={(checked, event) => {

                      let enable = checked ? u8TxSwitchEnable | (1 << i) : u8TxSwitchEnable & (~(1 << i))
                      setU8TxSwitchEnable(enable)
                      console.log("🚀 ~ file: pointset.tsx ~ line 300 ~ {[1,2,3,4,5,6].map ~ enable", enable)
                      // console.log("🚀 ~ file: profile.tsx ~ line 369 ~ ProfileView ~ checked,event",i, checked,event)

                    }}
                    checked={(u8TxSwitchEnable & (1 << i)) === (1 << i)} /></Col>
                </Row>
              </Col>)
            })}
          </Row>
        </Row>

      </Col>
    </Row>


  </>);
}

export default PointSet;