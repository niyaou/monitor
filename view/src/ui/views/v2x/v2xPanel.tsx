import React, { useState, useEffect } from 'react';
import { Row, Col, Switch, Form, Input, Divider, Button, Radio, message, Typography, Tooltip, notification, Checkbox, InputNumber, Space } from 'antd';
import { PageContainer } from '@ant-design/pro-layout';
import { SettingOutlined, CaretRightOutlined, ForwardOutlined, BugOutlined, EllipsisOutlined, FolderOpenOutlined, SaveOutlined } from '@ant-design/icons';
import type { CheckboxValueType } from 'antd/es/checkbox/Group';
import { find, filter } from 'lodash'
import baseApi from '../../api/fc2/baseApi'
import type { RadioChangeEvent } from 'antd';

/** 
 *  
*/

export interface Device {
    Kind :    string,
    Name   :  string,
    Protocol: string,
    Address : string,
}


function V2xPanel() {

    const commandsCAN = [{ group: 'CAN数据诊断', name: 'CAN数据丢失', id: 0x03000001, input: false }, { group: 'CAN数据诊断', name: 'CAN数据错误', id: 0x03000002, input: false }]

    const commandsGNSS = [{ group: 'GNSS数据诊断', name: '主车GNSS数据丢失', id: 0x04000001, input: false }, { group: 'GNSS数据诊断', name: '主车GNSS数据错误', id: 0x04000002, input: false }]

    const commandsV2x = [{ group: 'V2X数据诊断', name: 'RV数据丢失', id: 0x05000001, input: true }, { group: 'V2X数据诊断', name: 'MAP数据丢失', id: 0x05000003, input: true },
    { group: 'V2X数据诊断', name: 'SPAT数据丢失', id: 0x05000005, input: true },
    { group: 'V2X数据诊断', name: 'RSI数据丢失', id: 0x05000007, input: true }, { group: 'V2X数据诊断', name: 'RSM数据丢失', id: 0x05000009, input: true },
    { group: 'V2X数据诊断', name: 'BSM位置数据非法', id: 0x0500000c, input: true }, { group: 'V2X数据诊断', name: 'BSM编码失败', id: 0x0500000d, input: true },
    { group: 'V2X数据诊断', name: 'BSM_vehicleID非法', id: 0x0500000e, input: true },
    { group: 'V2X数据诊断', name: 'BSM_档位信息非法', id: 0x0500000f, input: true }, { group: 'V2X数据诊断', name: 'BSM_速度信息非法', id: 0x05000010, input: true },
    { group: 'V2X数据诊断', name: 'BSM_航向角信息非法', id: 0x05000011, input: true },
    { group: 'V2X数据诊断', name: 'BSM_转向角信息非法', id: 0x05000012, input: true }, { group: 'V2X数据诊断', name: 'BSM_加速度信息非法', id: 0x05000013, input: true },
    { group: 'V2X数据诊断', name: 'BSM_车身尺寸非法', id: 0x05000014, input: true },
    { group: 'V2X数据诊断', name: 'BSM_车辆类型非法', id: 0x05000015, input: true },
    { group: 'V2X数据诊断', name: 'V2X安全验签失败', id: 0x05000018, input: false },
    { group: 'V2X数据诊断', name: 'RSC数据丢失', id: 0x0500001a, input: false },
    { group: 'V2X数据诊断', name: 'VIR数据丢失', id: 0x0500001c, input: false },
    { group: 'V2X数据诊断', name: 'SSM数据丢失', id: 0x0500001e, input: false },
    { group: 'V2X数据诊断', name: 'VIR数据发送失败', id: 0x0500001f, input: false }]


    const commandsInner = [{ group: '内部处理错误', name: '主车CAN数据与GNSS数据时间同步失败', id: 0x08000001, input: false },
    { group: '内部处理错误', name: 'CAN&WSM消息类型非法', id: 0x08000003, input: false }, { group: '内部处理错误', name: 'RV与主车时间不同步', id: 0x08000004, input: false },
    { group: '内部处理错误', name: '车身数据不可识别的消息类型', id: 0x08000005, input: true }, { group: 'V2XService未拿到HV数据', name: '内部处理错误', id: 0x0800000a, input: false },
    { group: '内部处理错误', name: 'V2XService未拿到RV数据', id: 0x0800000b, input: false }, { group: '内部处理错误', name: 'V2XService未拿到RSU数据', id: 0x0800000c, input: false },]


    const [keyboard, setKeyboard] = useState();

    const [selectDevice, setSelectDevice] = useState('');


    const [devices, setDevices] = useState<Device[]>([{
        Kind :  '  string',
        Name   : ' string1',
        Protocol: 'string',
        Address :'192.168.0.0',
    },{
        Kind :  '  string',
        Name   : ' string333',
        Protocol: 'string',
        Address :'192.168.0.1',
    }]);

    const [inputValue, setInputValue] = useState([]);



    const [cankCmds, setCankCmds] = useState([]);
    const [gnssCmds, setGnssCmds] = useState([]);
    const [v2xCmds, setV2xCmds] = useState([]);
    const [innerCmds, setInnerCmds] = useState([]);

    const [mergedCmd, setMergedCmd] = useState([]);

    useEffect(() => {
        let selected = [...cankCmds, ...gnssCmds, ...v2xCmds, ...innerCmds]
        selected = selected.map(id => {
            let _param = find(inputValue, { id })
            if (typeof _param === 'undefined') {
                return { id, value: 0x999999 }
            } else {
                return { id, value: _param.value }
            }

        })
        console.log("🚀 ~ file: v2xPanel.tsx ~ line 59 ~ useEffect ~ selected", selected)
        setMergedCmd(selected)
    }, [inputValue, cankCmds, gnssCmds, v2xCmds, innerCmds])


    useEffect(() => {
        onInitDevice()
    }, [])






    const onInitDevice = () => {
        
        baseApi.getReq('device/list').then((res)=>{
            console.log("🚀 ~ file: v2xPanel.tsx ~ line 93 ~ baseApi.postReq ~ res", res)
            setDevices([])
            }).catch(err=>{
            console.log("🚀 ~ file: v2xPanel.tsx ~ line 95 ~ baseApi.postReq ~ err", err)
                
            })
    }



    const onChangeInputValue = (id, value) => {
        let _value = [...inputValue]

        let target = find(_value, { id })
        if (typeof target === 'undefined') {
            target = { id: id, value: value }
            _value.push(target)
        } else {
            target.value = value
        }

        console.log("🚀 ~ file: v2xPanel.tsx ~ line 50 ~ onChangeInputValue ~ target", target, _value)
        setInputValue(_value)
    }

const onVerifiedStart = () => {
   
        let device = find(devices, { Address:selectDevice })
        let params = {Dev:device,Event:mergedCmd }
        console.log("🚀 ~ file: v2xPanel.tsx ~ line 195 ~ V2xPanel ~ mergedCmd",params, mergedCmd,selectDevice)
        baseApi.postReq('vbox/event/v2xdiagnose',params).then((res)=>{
        console.log("🚀 ~ file: v2xPanel.tsx ~ line 93 ~ baseApi.postReq ~ res", res)

        }).catch(err=>{
        console.log("🚀 ~ file: v2xPanel.tsx ~ line 95 ~ baseApi.postReq ~ err", err)
            
        })
   
}

    const onCanChange = (checkedValues: CheckboxValueType[]) => {
        console.log('checked = ', checkedValues);
        setCankCmds(checkedValues)
    };
    const onGnssChange = (checkedValues: CheckboxValueType[]) => {
        console.log('checked = ', checkedValues);
        setGnssCmds(checkedValues)
    };
    const onV2xChange = (checkedValues: CheckboxValueType[]) => {
        console.log('checked = ', checkedValues);
        setV2xCmds(checkedValues)
    };
    const onInnerChange = (checkedValues: CheckboxValueType[]) => {
        console.log('checked = ', checkedValues);
        setInnerCmds(checkedValues)
    };



    const inputCommand = (cmds, params: any) => {
        return (
            <Col span={12}>
                <Row justify="start">
                    <Col span={12} style={{ height: '36px' }}>
                        <Checkbox
                            value={params.id}
                            onChange={(v) => {
                            }}
                            checked={cmds.indexOf(params.id) > -1}
                        >
                            {params.name}
                        </Checkbox>
                    </Col>
                    <Col span={12}>

                        {params.input && (<InputNumber min={1} max={10} disabled={cmds.indexOf(params.id) < 0} keyboard={keyboard}
                            onChange={(value => {
                                onChangeInputValue(params.id, value)
                            })}
                        />)}
                    </Col>
                </Row>
            </Col>)
    }


    return (
        <PageContainer
            title='V2X测试面板'
            // subTitle={resultStatusComp}
            extra={
                [
                    <Button key="extra-load1" disabled={mergedCmd.length === 0 ||selectDevice===''} icon={<CaretRightOutlined />}
                        onClick={onVerifiedStart}
                    >开始测试</Button>,
                ]}
        >
             <Divider orientation="left" orientationMargin="0">在线设备</Divider>
             <Row gutter={24}>
             <Radio.Group onChange={(v)=>{
             console.log("🚀 ~ file: v2xPanel.tsx ~ line 206 ~ V2xPanel ~ v", v)
             setSelectDevice(v.target.value)
             }} value={selectDevice}>
             {devices.map((device,i) =>(<Radio key={`radio-${i}`} value={device.Address}>{device.Name}</Radio>))     }
      
    
    </Radio.Group>
                </Row>
            <Divider orientation="left" orientationMargin="0">CAN数据诊断</Divider>
            <Row gutter={24}>
                {/* <Col span={24}>
                    <Title level={5}>h4. Ant Design</Title>
                </Col> */}
                <Col span={24}>

                    <Checkbox.Group style={{ width: '100%' }} onChange={onCanChange}>
                        <Row gutter={24}>
                            {commandsCAN.map((command) => inputCommand(cankCmds, command))}
                        </Row>
                    </Checkbox.Group>

                </Col>
            </Row>
            <Divider orientation="left" orientationMargin="0">GNSS数据诊断</Divider>
            <Row gutter={24}>

                <Col span={24}>
                    <Checkbox.Group style={{ width: '100%' }} onChange={onGnssChange}>
                        <Row gutter={24}>
                            {commandsGNSS.map((command) => inputCommand(gnssCmds, command))}
                        </Row>
                    </Checkbox.Group>
                </Col>
            </Row>
            <Divider orientation="left" orientationMargin="0">V2X数据诊断</Divider>
            <Row gutter={24}>

                <Col span={24}>
                    <Checkbox.Group style={{ width: '100%' }} onChange={onV2xChange}>
                        <Row gutter={24}>
                            {commandsV2x.map((command) => inputCommand(v2xCmds, command))}
                        </Row>
                    </Checkbox.Group>
                </Col>
            </Row>
            <Divider orientation="left" orientationMargin="0">内部处理错误</Divider>
            <Row gutter={24}>

                <Col span={24}>
                    <Checkbox.Group style={{ width: '100%' }} onChange={onInnerChange}>
                        <Row gutter={24}>
                            {commandsInner.map((command) => inputCommand(innerCmds, command))}
                        </Row>
                    </Checkbox.Group>
                </Col>
            </Row>
        </PageContainer >
    );
}

export default V2xPanel;
