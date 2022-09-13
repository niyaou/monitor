import React, { useState, useEffect } from 'react';
import { Row, Col, Tabs, Form, Input, Switch, Radio, Checkbox, Select, InputNumber, FormInstance } from 'antd';
import { IForm } from './types'
import lodash from 'lodash'

const { Option } = Select;

const fieldMap = {
  "u32RxChannelEnable": "strRfeSetting.u32RxChannelEnable",
}

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



const Common: React.FC<IForm> = ({ data, formName, isSubmitStatus, setFailed, setValue, directlySetValue, registForm }: IForm) => {
  const [form] = Form.useForm();

  const [arange, setArange] = useState<number>(lodash.get(data, "strTef82xxFrameOpParam.u8CafcPllLPFLUTSel"));

  // const defaultValue = {
  //   "strRfeSetting.u8FrontEnd": lodash.get(data, "strRfeSetting.u8FrontEnd"),
  //   "u16AcqNrFrames": lodash.get(data, "u16AcqNrFrames"),
  //   "strTef82xxFrameOpParam.u8ProfStayCnt": lodash.get(data, "strTef82xxFrameOpParam.u8ProfStayCnt"),
  //   "strTef82xxFrameOpParam.u8CafcPllLPFLUTSel": lodash.get(data, "strTef82xxFrameOpParam.u8CafcPllLPFLUTSel"),
  //   "strTef82xxFrameOpParam.u32CafcLoopBandwidth": lodash.get(data, "strTef82xxFrameOpParam.u32CafcLoopBandwidth"),
  //   "strRfeSetting.u32SamplingFrequency": lodash.get(data, "strRfeSetting.u32SamplingFrequency"),
  //   "strRfeSetting.u16NrChirpsInFrame": lodash.get(data, "strRfeSetting.u16NrChirpsInFrame"),
  //   "strRfeSetting.u16NrSamplesPerChirp": lodash.get(data, "strRfeSetting.u16NrSamplesPerChirp"),
  //   "strRfeSetting.u8NrChirpShapes": lodash.get(data, "strRfeSetting.u8NrChirpShapes"),
  // }
  // const u32RxChannelEnable = lodash.get(data, fieldMap["u32RxChannelEnable"])
  // for (let index = 0; index < 8; index++) {
  //   defaultValue[`rxChannelEnable[${index}]`] = Boolean(getBit(u32RxChannelEnable, index))
  // }
  const [defaultValue, setDefaultValue] = useState({});


  useEffect(() => {
    // console.log('common data change:', data)
    const _defaultValue = {
      "strRfeSetting.u8FrontEnd": lodash.get(data, "strRfeSetting.u8FrontEnd"),
      "u16AcqNrFrames": lodash.get(data, "u16AcqNrFrames"),
      "strTef82xxFrameOpParam.u8ProfStayCnt": lodash.get(data, "strTef82xxFrameOpParam.u8ProfStayCnt"),
      "strTef82xxFrameOpParam.u8CafcPllLPFLUTSel": lodash.get(data, "strTef82xxFrameOpParam.u8CafcPllLPFLUTSel"),
      "strTef82xxFrameOpParam.u32CafcLoopBandwidth": lodash.get(data, "strTef82xxFrameOpParam.u32CafcLoopBandwidth"),
      "strRfeSetting.u32SamplingFrequency": lodash.get(data, "strRfeSetting.u32SamplingFrequency"),
      "strRfeSetting.u16NrChirpsInFrame": lodash.get(data, "strRfeSetting.u16NrChirpsInFrame"),
      "strRfeSetting.u16NrSamplesPerChirp": lodash.get(data, "strRfeSetting.u16NrSamplesPerChirp"),
      "strRfeSetting.u8NrChirpShapes": lodash.get(data, "strRfeSetting.u8NrChirpShapes"),
    }
    const u32RxChannelEnable = lodash.get(data, fieldMap["u32RxChannelEnable"])
    for (let index = 0; index < 8; index++) {
      _defaultValue[`rxChannelEnable[${index}]`] = Boolean(getBit(u32RxChannelEnable, index))
    }
    setDefaultValue(_defaultValue)
  }, [data]);

  useEffect(() => {
    // console.log(formName, isSubmitStatus)
    if (isSubmitStatus) {
      form.submit()
    }
  }, [isSubmitStatus]);


  useEffect(() => {
    form.resetFields()
  }, [defaultValue]);

  // useEffect(() => {
  //   console.log("🚀 ~ file: common.tsx ~ line 61 ~ arange", arange)
  // },[arange])

  const onFinish = (values: any) => {
    // console.log('Success:', values);

    const reformValue = getReformValue(values);
    setValue(formName, reformValue);
  };

  const onLUTChange = (value: number) => {
    setArange(value)

  };

  const onFinishFailed = (errorInfo: any) => {
    console.log('Failed:', errorInfo);
    setFailed()
  };

  const onValuesChange = (changedValues, allValues) => {
    console.log(changedValues, allValues);
    directlySetValue(getReformValue(allValues))
  };

  useEffect(() => {
    registForm(formName)

  }, []);

  // useEffect(() => {
  //   console.log(' [data]',data)

  // }, [data]);


  const rxChannelEnables = new Array(8)
  for (let index = 0; index < 8; index++) {
    rxChannelEnables[index] = (
      <Form.Item key={`rxChannelEnable${index}`}
        label={`接收天线使能_${index + 1}`}
        valuePropName="checked"
        name={`rxChannelEnable[${index}]`}
      >
        <Switch checkedChildren="使能" unCheckedChildren="失能" />
      </Form.Item>
    )
  }

  return <>
    <Form
      form={form}
      name={formName}
      labelCol={{ span: 9 }}
      wrapperCol={{ span: 15 }}
      initialValues={defaultValue}
      onFinish={onFinish}
      onFinishFailed={onFinishFailed}
      autoComplete="off"
      onValuesChange={onValuesChange}
    >
      <Row gutter={24}>
        <Col offset={6} span={6}>
          <Form.Item
            label="主从模式"
            name="strRfeSetting.u8FrontEnd"
          >
            <Select>
              <Option value={1}>master</Option>
              <Option value={2}>master and slaver</Option>
            </Select>
          </Form.Item>
          <Form.Item
            label="采帧次数"
            name="u16AcqNrFrames"
            tooltip=" 范围:[1,65535]"
            rules={[{ required: true }]}
          >
            <InputNumber style={{ width: '100%' }} min={1} max={65535} />
          </Form.Item>
          <Form.Item
            label="Profile重复次数"
            name="strTef82xxFrameOpParam.u8ProfStayCnt"
            tooltip=" 范围:[0,255]"
            rules={[{
              required: true,
              // type: 'number', min: 0, max: 65535 
            }]}
          >
            <InputNumber style={{ width: '100%' }} min={1} max={255} />
          </Form.Item>
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
              onChange={onLUTChange}
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
            tooltip={`范围:[${arange === 0 ? 200000 : arange === 1 ? 250000 : 300000} ~ ${arange === 0 ? 1600000 : arange === 1 ? 1650000 : 1700000}]`}
            rules={[{
              required: true,
              // type: 'number', min: 0, max: 65535 
            }]}
          >
            {/* <Select
              dropdownStyle={{ background: "rbg(240,242,245)" }}
              style={{ width: '100%', background: '#00000000' }}
              bordered={false}
            >
              {['20000 ~ 1600000',
                '25000 ~ 1650000',
                '30000 ~ 1700000'
              ].map((v, i) => { return (<Option key={`select-${i}`} value={i} style={{ background: '#00000000' }}>{v}</Option>) })}
            </Select> */}

            <InputNumber style={{ width: '100%' }} min={arange === 0 ? 200000 : arange === 1 ? 250000 : 300000} max={arange === 0 ? 1600000 : arange === 1 ? 1650000 : 1700000} />
          </Form.Item>

          <Form.Item
            label="采样频率"
            name="strRfeSetting.u32SamplingFrequency"
            tooltip=" 范围:(5000 ,40000 )"
            rules={[{ required: true }]}
          >
            <InputNumber style={{ width: '100%' }} min={5000} max={40000} />
          </Form.Item>
          <Form.Item
            label="chirp 数"
            name="strRfeSetting.u16NrChirpsInFrame"
            tooltip="一帧里面chirp数量,范围:0 ~ 4095"
            rules={[{ required: true }]}
          >
            <InputNumber style={{ width: '100%' }} min={0} max={4095} />
          </Form.Item>
          <Form.Item
            label="每个chirp采样点数"
            name="strRfeSetting.u16NrSamplesPerChirp"
            tooltip="每个chirp 采样点数,范围：(0,4095)"
            rules={[{ required: true }]}
          >
            <InputNumber style={{ width: '100%' }} min={0} max={4095} />
          </Form.Item>
          <Form.Item
            label="Chirp配置表数量"
            name="strRfeSetting.u8NrChirpShapes"
            tooltip="范围:1~8"
            rules={[{ required: true }]}
          >
            <InputNumber style={{ width: '100%' }} min={1} max={8} />
          </Form.Item>
        </Col>
        <Col span={12}>
          {rxChannelEnables}
        </Col>
      </Row>

    </Form>
  </>
};

export default Common;

function getReformValue(values: any) {
  const reformValue = {};
  Object.keys(values).forEach(key => {
    if (!key.startsWith("rxChannelEnable")) {
      lodash.set(reformValue, key, values[key]);
    }
  });
  let rxChannelEnableValue = 0;
  for (let index = 0; index < 8; index++) {
    if (values[`rxChannelEnable[${index}]`] === true) {
      rxChannelEnableValue = rxChannelEnableValue | 1 << index;
    }
  }
  lodash.set(reformValue, fieldMap["u32RxChannelEnable"], rxChannelEnableValue);
  return reformValue;
}

