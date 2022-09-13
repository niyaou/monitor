import React, { useState, useEffect, useContext, useRef } from 'react';
import { Row, Col, Card, InputNumber, Form, Spin, Switch, Statistic, Button, Alert, PageHeader, notification } from 'antd';
import { DownloadOutlined, UploadOutlined, PlaySquareOutlined, StopOutlined, EllipsisOutlined, FolderOpenOutlined, SaveOutlined } from '@ant-design/icons';
import turretService from '../../../service/turretService'
import { configDataContext, actionTypes, sessionFlag } from "./configDataReducer";
import { status } from '@grpc/grpc-js';
import { cloneDeep, find, findIndex, drop } from 'lodash'
import configApi from '../../../api/fc2/configApi'
import { parseDataByProtocal } from './config2'
import baseApi from '../../../api/fc2/baseApi';
const defaultValue: object = {


    "Y": {
        "range": 0,
        "maxSpeed": 0
    },
    "X": {
        "range": 0,
        "maxSpeed": 0
    }
}

const defaultData: object = {
    "X": 0, "Y": 0
}

const TurretWorkStatus: React.FC<any> = (props) => {
    const { workStatus } = props
    const getAlertType = (workStatus) => {
        switch (workStatus) {
            case 'work':
                return 'success'
            case 'idle':
                return 'info'
            default:
                return 'warning'
        }
    }

    const alertType = getAlertType(workStatus)
    return (<Alert message={workStatus} type={alertType} showIcon />)
}

const TurretStatus: React.FC<any> = (props) => {
    const { data } = props

    return (
        <Card title="运行状态" size='small'>
            {/* <Spin spinning={statusLoading}> */}
            <Row gutter={24}>

                <Col span={8}>
                    <Statistic title="俯仰角" value={data["Y"]} />
                </Col>
                <Col span={8}>
                    <Statistic title="方位角" value={data["X"]} />
                </Col>
            </Row>
            {/* </Spin> */}
        </Card>
    )
}




const AngleConfig: React.FC = () => {
    const { state, dispatch } = useContext(configDataContext);
    const [angleForm] = Form.useForm();
    const [startLoading, setStartLoading] = useState(false);
    const [partAngleName, setPartAngleName] = useState('');
    const [setAngleLoading, setSetAngleLoading] = useState(false);

    const setAngle = (values) => {
        turretService.setAngle(values)
            .then(function (data) {
                notification.success({
                    message: `操作成功`,
                    // description: '设置角度范围成功。',
                    placement: 'topRight',
                    duration: 1,
                });
                // setData(data)
                // setWorkStatus(data.workStatus)
            })
            .catch(function (error) {
                console.error(error);
                notification.error({
                    message: `操作失败`,
                    description: `${error}`,
                });
            })
            .finally(() => { setSetAngleLoading(false) })
    };

    const setPartAngle = (values) => {
        turretService.setPartAngle(values)
            .then(function (data) {
                notification.success({
                    message: `操作成功`,
                    placement: 'topRight',
                    duration: 1,
                });
            })
            .catch(function (error) {
                console.error(error);
                notification.error({
                    message: `操作失败`,
                    description: `${error}`,
                });
            })
            .finally(() => { setPartAngleName(''); setSetAngleLoading(false) })
    };

    const onFinish = (values: any) => {
        if (setAngleLoading) {
            setAngle(values)
        } else if (partAngleName !== '') {
            setPartAngle({
                name: partAngleName,
                value: values[partAngleName]['range']
            })
        } else {
            console.warn("no next step!");
        }
    };

    return (
        <Card title="运行角度" size='small' extra={[
            <Button
                key="extra-play"
                icon={<PlaySquareOutlined />}
                loading={startLoading}
                disabled={state.flag === sessionFlag.FLAG_START_SESSION}
                onClick={() => { setSetAngleLoading(true); angleForm.submit() }}
            >
                设置转台
            </Button>
        ]}>
            <Form name="angle-form"
                form={angleForm}
                labelCol={{ span: 8 }}
                wrapperCol={{ span: 16 }}
                onFinish={onFinish}
                initialValues={defaultValue}
            >

                <Row gutter={24}>
                    <Col span={12}>
                        <Form.Item
                            label="俯仰角"
                            name={['Y', 'range']}
                            rules={[{ required: true }]}
                            tooltip=" 范围:±44°"
                        >
                            <InputNumber min={-44} max={44} addonAfter="°" />
                        </Form.Item>
                    </Col>


                    <Col span={12}>
                        <Form.Item
                            label="方位角"
                            name={['X', 'range']}
                            rules={[{ required: true }]}
                            tooltip=" 范围:±180°"
                        >
                            <InputNumber min={-180} max={180} addonAfter="°" />
                        </Form.Item>
                    </Col>

                </Row>
            </Form>
        </Card>
    )
}





const SessionConfig: React.FC<any> = (props) => {
    const [sessionForm] = Form.useForm();
    const [startLoading, setStartLoading] = useState(false);
    const [partAngleName, setPartAngleName] = useState('');
    const [setAngleLoading, setSetAngleLoading] = useState(false);
    const [angleOrPosition, setAngleOrPosition] = useState(false);

    const [sessionParams, setSessionParams] = useState({ duration: 1, Steps: [], waitTrigger: false });

    const [current, setCurent] = useState({ Y: 0, X: 0 });

    const { state, dispatch } = useContext(configDataContext);
    // const { data: { Polar, EL, AZ }, workStatus } = props
    const { data } = props
    const sessionParamsRef = useRef()
    const sendToRadar = (configData) => {
        try {
            notification.info({
                message: `发出连续采集设置`,
                placement: 'topRight',
                duration: 3,
            });
            configApi.set(parseDataByProtocal(configData)).then((response) => {
                const data = response.data
                if (data.enuParamCfgAck === 1) {
                    notification.success({
                        message: `雷达设置完成`,
                        placement: 'topRight',
                        duration: 3,
                    });
                } else {
                    console.error('错误:' + data.enuParamCfgAck);
                    notification.error({
                        message: `雷达操作失败1`,
                        duration: 3,
                        description: `错误内容： ${data.enuParamCfgAck}    ${response.msg}`,
                    });
                }
            }).catch((error) => {
                console.error(error);
                notification.error({
                    message: `雷达操作失败2`,
                    duration: 3,
                    description: `错误内容： ${error.message}`,
                });
            })
        } catch (error) {
            console.error(error);
            notification.error({
                message: `雷达操作失败3`,
                duration: 3,
                description: `错误内容：  ${error.message}`,
            });
        }
    }


    const initialSessionParam = (sessionForm: any) => {
        //解析连续发波参数，初始化数据
        let param = cloneDeep(sessionParams)


        let _azStep = []
        if (sessionForm.AZStep.range > 0) {
            for (let i = sessionForm.XMin.range; i <= sessionForm.XMax.range; i += sessionForm.XStep.range) {
                _azStep.push(i)
            }
        }



        let _elStep = []
        if (sessionForm.YStep.range > 0) {
            for (let i = sessionForm.YMin.range; i <= sessionForm.YMax.range; i += sessionForm.YStep.range) {
                _elStep.push(i)
            }
        } else {
            _elStep.push(sessionForm.YMin.range)
        }

        let step = _elStep.map((elStep) => { return [elStep, _azStep] })



        // param.azSteps = _azStep
        param.Steps = step
        param.duration = sessionForm.Duration.range
        console.log("🚀 ~ file: turret.tsx ~ line 438 ~ initialSessionParam ~ param", param)
        setSessionParams(param)
    }

    const startSequence = (sessionParams) => {
        //开始按照序列发波
        if (sessionParams.Steps.length === 0) {
            dispatch({ type: actionTypes.FC_SESSION, ...{ payload: { flag: sessionFlag.FLAG_IDLE_SESSION } } })
            return
        }

        if (sessionParams.waitTrigger) {
            return
        }
        // if(data.data.EL.range !== ){

        // }
        let values = { Y: { range: 0 }, X: { range: 0 } }

        let turrentParams = sessionParams.Steps[0]

        values.Y.range = turrentParams[0]
        values.X.range = turrentParams[1][0]

        setCurent({ Y: values.Y.range, X: values.X.range })

        let _params = cloneDeep(sessionParams)
        _params.waitTrigger = true
        setSessionParams(_params)
        let _state = cloneDeep(state)
        _state.u16Cmd = 2
        console.log("🚀 ~ file: turret.tsx ~ line 398 ~ startSequence ~ values", values, _state)
        sendToRadar(_state)
        setTimeout(() => {
            // let _params = cloneDeep(sessionParams)
            // setSessionParams(_params)

            let _params = cloneDeep(sessionParams)
            if (_params.Steps.length > 0) {
                if (_params.Steps[0][1].length === 1) {
                    _params.Steps = drop(_params.Steps)
                } else {
                    _params.Steps[0][1] = drop(_params.Steps[0][1])
                }

            } else {
                _params.Steps = drop(_params.Steps)
            }

            _params.waitTrigger = false
            setSessionParams(_params)
            turretService.setAngle(values).then(() => {
                notification.success({
                    message: `转台设置成功${JSON.stringify(values)}`,
                    // description: '设置角度范围成功。',
                    placement: 'topRight',
                    duration: 3,
                });
                setTimeout(() => {
                    turretService.start()
                        .then(function (data) {
                            notification.success({
                                message: `转台启动成功${JSON.stringify(values)}`,
                                // description: '设置角度范围成功。',
                                placement: 'topRight',
                                duration: 5,
                            });

                        })
                }, 600)

            })
                .catch(function (error) {
                    console.error(error);
                    notification.error({
                        message: `转台操作失败`,
                        description: `${error}`,
                    });
                })
                .finally(() => { setSetAngleLoading(false) })

        }, sessionParams.duration * 1000)

    }

    useEffect(() => {
        if (data.workStatus === 'idle') {
            startSequence(sessionParams)
        }
    }, [data])

    // useEffect(() => {
    //     console.log("🚀 ~ file: turret.tsx ~ line 394 ~ sessionParams", sessionParams)
    // }, [sessionParams])


    const getAlertType = (workStatus) => {
        switch (workStatus) {
            case 'work':
                return 'success'
            case 'idle':
                return 'info'
            default:
                return 'warning'
        }
    }

    return (
        <Card title="连续采集" size='small' extra={[
            <Button
                key="extra-play"
                icon={<PlaySquareOutlined />}
                loading={startLoading}
                disabled={state.flag === sessionFlag.FLAG_START_SESSION}
                onClick={() => {
                    setSetAngleLoading(true);
                    sessionForm.submit()


                }}
            >
                开始连续采集
            </Button>
        ]}>
            <Form name="angle-form"
                form={sessionForm}
                labelCol={{ span: 8 }}
                wrapperCol={{ span: 16 }}
                onFinish={(v) => {

                    console.log("🚀 ~ file: turret.tsx ~1 line 380 ~ v", v)
                    initialSessionParam(v)
                    dispatch({ type: actionTypes.FC_SESSION, ...{ payload: { flag: sessionFlag.FLAG_START_SESSION } } })

                }}
                initialValues={defaultValue}
            >
                <Row gutter={1} justify="start">
                    <Col span={6}>
                        <Form.Item
                            label={'俯仰低'}
                            name={['YMin', 'range']}
                            rules={[{ required: true, }]}
                            tooltip={'范围:±180°'}
                        >
                            <InputNumber min={-180} max={180} addonAfter="°" />
                        </Form.Item>
                    </Col>
                    <Col span={6}>
                        <Form.Item
                            label={'俯仰高'}
                            name={['YMax', 'range']}
                            rules={[{ required: true, }]}
                            tooltip={'范围:±180°'}
                        >
                            <InputNumber min={-180} max={180} addonAfter="°" />
                        </Form.Item>
                    </Col>
                    <Col span={6}>
                        <Form.Item
                            label={'俯仰步长'}
                            name={['YStep', 'range']}
                            rules={[{ required: true, }]}
                            tooltip={'范围:0~90°'}
                        >
                            <InputNumber min={0} max={90} addonAfter="°" />
                        </Form.Item>
                    </Col>

                </Row>
                <Row gutter={1} justify="start">



                    <Col span={6}>
                        <Form.Item
                            label={'方位低'}
                            name={['XMin', 'range']}
                            rules={[{ required: true }]}
                            tooltip={'范围:±180°'}
                        >
                            <InputNumber min={-180} max={180} addonAfter="°" />
                        </Form.Item>
                    </Col>
                    <Col span={6}>
                        <Form.Item
                            label={'方位高'}
                            name={['XMax', 'range']}
                            rules={[{ required: true }]}
                            tooltip={'范围:±180°'}
                        >
                            <InputNumber min={-180} max={180} addonAfter="°" />
                        </Form.Item>
                    </Col>
                    <Col span={6}>
                        <Form.Item
                            label={'方位步长'}
                            name={['XStep', 'range']}
                            rules={[{ required: true, }]}
                            tooltip={'范围:0~90°'}
                        >
                            <InputNumber min={0} max={90} addonAfter="°" />
                        </Form.Item>
                    </Col>
                    <Col span={6}>
                        <Form.Item
                            label="采集时间"
                            name={['Duration', 'range']}
                            rules={[{ required: true }]}
                            tooltip="1~120秒"
                        >
                            <InputNumber min={1} max={120} addonAfter="s" />
                        </Form.Item>
                    </Col>
                </Row>


            </Form>
        </Card>
    )
}


const TurretTable: React.FC = () => {
    const { state, dispatch } = useContext(configDataContext);
    // const [workStatus, setWorkStatus] = useState('unkonw');
    const [data, setData] = useState(defaultData);
    const [startLoading, setStartLoading] = useState(false);
    const [endLoading, setEndLoading] = useState(false);
    useEffect(() => {
        getStatus()
        const timeoutID = setInterval(
            getStatus
            , 500000);
        return () => {
            clearInterval(timeoutID)
        }
    }, []);

    useEffect(() => {
        if (endLoading === false)
            return
        stop()
    }, [endLoading])

    useEffect(() => {
        if (startLoading === false)
            return
        turretService.start()
            .then(function (data) {
                notification.success({
                    message: `操作成功`,
                    // description: '开始运行成功。',
                    placement: 'topRight',
                    duration: 1,
                });
                // setData(data)
                // setWorkStatus(data.workStatus)
            })
            .catch(function (error) {
                console.error(error);
                notification.error({
                    message: `操作失败1`,
                    description: `${error}`,
                });
            })
            .finally(() => { setStartLoading(false) })
    }, [startLoading]);

    const stop = () => {
        turretService.stop()
            .then(function (data) {
                notification.success({
                    message: `操作成功`,
                    // description: '停止运行成功。',
                    placement: 'topRight',
                    duration: 1,
                });
                // setData(data)
                // setWorkStatus(data.workStatus)
            })
            .catch(function (error) {
                console.error(error);
                notification.error({
                    message: `操作失败2`,
                    description: `${error}`,
                });
            })
            .finally(() => { setEndLoading(false) })
    };

    const getStatus = () => {

        baseApi.getReq('RotartTableCommandAcknowledge').then(function (data) {
            setData(data)
            // setWorkStatus(data.workStatus)
        }).catch(function (error) {
            console.error(error);
            notification.error({
                message: `操作失败3`,
                description: `${error}`,
            });
        })
    };

    return (
        <>
            <PageHeader
                title='转台'


            />
            <Row gutter={24} style={{ marginTop: '8px' }}>
                <Col offset={2} span={6}>
                    {/* <LaserSwitch /> */}
                </Col>
            </Row>
            <Row gutter={24}>
                <Col offset={2} span={16}>
                    <AngleConfig />
                </Col>
                <Col span={8}>
                    {/* <SpeedConfig /> */}
                </Col>
            </Row>
            <Row gutter={24}>
                <Col offset={2} span={16}>
                    <hr></hr>
                </Col>
                <Col offset={2} span={16}>
                    <TurretStatus data={data} />
                </Col>
            </Row>
            <Row gutter={24}>
                <Col offset={2} span={16}>
                    <hr></hr>
                </Col>
                <Col offset={2} span={16}>
                    <SessionConfig data={data} />
                </Col>
            </Row>
        </>
    )
};

export default TurretTable;

