import React, { useState, useEffect } from 'react';
import { Row, Col, Tabs, Form, Input, Button, Switch, Checkbox, Select, InputNumber, FormInstance } from 'antd';
import { IForm } from './types'
import lodash from 'lodash'

const fieldMap = {
  "u8aUseDDMA": "strTef82xxFrameOpParam.strPhaseRotator.u8aUseDDMA",
  "u32aDdmaInitPhase": "strTef82xxFrameOpParam.strPhaseRotator.u32aDdmaInitPhase",
  "u32aDdmaPhaseUpdate": "strTef82xxFrameOpParam.strPhaseRotator.u32aDdmaPhaseUpdate",
}

const DDMA: React.FC<IForm> = ({ data, formName, isSubmitStatus, setFailed, setValue, directlySetValue, registForm }: IForm) => {
  const [form] = Form.useForm();

  const defaultValue: object = {}
  Object.keys(fieldMap).forEach(key => {
    for (let index = 0; index < 6; index++) {
      if (key === "u8aUseDDMA") {
        defaultValue[`${key}[${index}]`] = Boolean(lodash.get(data, `${fieldMap[key]}[${index}]`))
      } else {
        defaultValue[`${key}[${index}]`] = lodash.get(data, `${fieldMap[key]}[${index}]`)
      }
    }
  });

  useEffect(() => {
    // console.log(formName, isSubmitStatus)


    if (isSubmitStatus) {
      form.submit()
    }
  }, [isSubmitStatus]);

  const onFinish = (values: any) => {
    // console.log('Success:', values);
    const reformValue = getReformValue(values);
    setValue(formName, reformValue);
  };

  const onValuesChange = (changedValues, allValues) => {
    console.log(changedValues, allValues);
    directlySetValue(getReformValue(allValues))
  };

  const onFinishFailed = (errorInfo: any) => {
    console.log('Failed:', errorInfo);
    setFailed()
  };

  useEffect(() => {
    // console.log(' [data]---------',data)
    form.resetFields()
  }, [data]);

  useEffect(() => {
    registForm(formName)
  }, []);

  const aUseDDMAs = new Array(6)
  const aDdmaInitPhases = new Array(6)
  const aDdmaPhaseUpdate = new Array(6)
  for (let index = 0; index < 6; index++) {
    // const profileKey = `profile_${index + 1}`
    aUseDDMAs[index] = (
      <Form.Item key={`u8aUseDDMA${index}`}
        label={`DDMA使能状态_${index + 1}`}
        valuePropName="checked"
        name={`u8aUseDDMA[${index}]`}
      >
        <Switch checkedChildren="使能" unCheckedChildren="失能" />
      </Form.Item>
    )
    aDdmaInitPhases[index] = (
      <Form.Item key={`u32aDdmaInitPhase${index}`}
        label={`DDMA初始相位_${index + 1}`}
        name={`u32aDdmaInitPhase[${index}]`}
        rules={[{ required: true, }]}
      >
        <InputNumber style={{ width: '100%' }} min={0} max={360.000} />
      </Form.Item>
    )
    aDdmaPhaseUpdate[index] = (
      <Form.Item key={`u32aDdmaPhaseUpdate${index}`}
        label={`DDMA更新相位_${index + 1}`}
        name={`u32aDdmaPhaseUpdate[${index}]`}
        rules={[{
          required: true,
          // type: 'number', min: 0, max: 65535 
        }]}
      >
        <InputNumber style={{ width: '100%' }} min={0} max={360.000} />
      </Form.Item>
    )
  }

  return <>
    <Form
      form={form}
      name="basic"
      labelCol={{ span: 8 }}
      wrapperCol={{ span: 16 }}
      initialValues={defaultValue}
      onFinish={onFinish}
      onFinishFailed={onFinishFailed}
      autoComplete="off"
      onValuesChange={onValuesChange}
    >
      <Row gutter={24}>
        <Col offset={2} span={6}>
          {aDdmaInitPhases}
        </Col>
        <Col span={6}>
          {aDdmaPhaseUpdate}
        </Col>
        <Col span={6}>
          {aUseDDMAs}
        </Col>
      </Row>
    </Form>
  </>
};

export default DDMA;

function getReformValue(values: any) {
  const reformValue = {};
  Object.keys(values).forEach(key => {
    // debugger
    const fieldName = key.substring(0, key.length - 3);
    const arrayIndex = key.substring(key.length - 3, key.length);
    if (fieldName === "u8aUseDDMA") {
      lodash.set(reformValue, fieldMap[fieldName] + arrayIndex, Number(values[key]));
    } else {
      lodash.set(reformValue, fieldMap[fieldName] + arrayIndex, values[key]);
    }
  });
  return reformValue;
}

