import React, { useState, useEffect } from 'react';
import { Row, Col, Card, InputNumber, Form, Spin, Switch, Statistic, Button, Alert, PageHeader, notification } from 'antd';
import { DownloadOutlined, UploadOutlined, PlaySquareOutlined, StopOutlined, EllipsisOutlined, FolderOpenOutlined, SaveOutlined } from '@ant-design/icons';
import darkroomModuleService from '../../../service/darkroomModuleService'

const defaultValue: object = {
    "line1": {
        "range": 0,
        "maxSpeed": 0
    },
    "line2": {
        "range": 0,
        "maxSpeed": 0
    },
    "radian1": {
        "range": 0,
        "maxSpeed": 0
    },
    "radian2": {
        "range": 0,
        "maxSpeed": 0
    },
    "dividedCircle": {
        "range": 0,
        "maxSpeed": 0
    }
}

const defaultData: object = {
    "workStatus": "unkonw",
    "data": {
        "line1": { "angle": 0, "speed": 0 },
        "line2": { "angle": 0, "speed": 0 },
        "radian1": { "angle": 0, "speed": 0 },
        "radian2": { "angle": 0, "speed": 0 },
        "dividedCircle": { "angle": 0, "speed": 0 },
    }
}

const WorkStatus: React.FC<any> = (props) => {
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

const ModuleStatus: React.FC<any> = (props) => {
    const { data } = props
    return (
        <Card title="运行状态" size='small'>
            {/* <Spin spinning={statusLoading}> */}
            <Row gutter={24}>
                <Col span={6}>
                    <Statistic title="直线模组1号运行位置(mm)" value={data["data"]["line1"]["angle"]} />
                </Col>
                <Col span={6}>
                    <Statistic title="直线模组2号运行位置(mm)" value={data["data"]["line2"]["angle"]} />
                </Col>
                <Col span={4}>
                    <Statistic title="弧形模组1号(度)" value={data["data"]["radian1"]["angle"]} />
                </Col>
                <Col span={4}>
                    <Statistic title="弧形模组2号(度)" value={data["data"]["radian2"]["angle"]} />
                </Col>
                <Col span={4}>
                    <Statistic title="分度盘(度)" value={data["data"]["dividedCircle"]["angle"]} />
                </Col>
            </Row>
            {/* </Spin> */}
        </Card>
    )
}

const AngleConfig: React.FC = () => {
    const [angleForm] = Form.useForm();
    const [startLoading, setStartLoading] = useState(false);
    const [partAngleName, setPartAngleName] = useState('');
    const [setAngleLoading, setSetAngleLoading] = useState(false);

    const setAngle = (values) => {
        darkroomModuleService.setAngle(values)
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
        darkroomModuleService.setPartAngle(values)
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
                onClick={() => { setSetAngleLoading(true); angleForm.submit() }}
            >
                设置全部
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
                    <Col span={18}>
                        <Form.Item
                            label="直线模组1号"
                            name={['line1', 'range']}
                            rules={[{ required: true, }]}
                            tooltip=" 范围:0-3750mm"
                        >
                            <InputNumber min={0} max={3750} addonAfter="mm" />
                        </Form.Item>
                    </Col>
                    <Col span={6}>
                        <Button
                            key="extra-play"
                            icon={<PlaySquareOutlined />}
                            loading={startLoading}
                            onClick={() => { setPartAngleName('line1'); angleForm.submit() }}
                        >
                            设置
                        </Button></Col>
                </Row>
                <Row gutter={24}>
                    <Col span={18}>
                        <Form.Item
                            label="直线模组2号"
                            name={['line2', 'range']}
                            rules={[{ required: true, }]}
                            tooltip=" 范围:0-3750mm"
                        >
                            <InputNumber min={0} max={3750} addonAfter="mm" />
                        </Form.Item>
                    </Col>
                    <Col span={6}>
                        <Button
                            key="extra-play"
                            icon={<PlaySquareOutlined />}
                            loading={startLoading}
                            onClick={() => { setPartAngleName('line2'); angleForm.submit() }}
                        >
                            设置
                        </Button></Col>
                </Row>
                <Row gutter={24}>
                    <Col span={18}>
                        <Form.Item
                            label="弧度模组1号"
                            name={['radian1', 'range']}
                            rules={[{ required: true, }]}
                            tooltip="-90°至【当前弧形模组2号位置减2度】"
                        >
                            <InputNumber min={-90} max={90} addonAfter="°" />
                        </Form.Item>
                    </Col>
                    <Col span={6}>
                        <Button
                            key="extra-play"
                            icon={<PlaySquareOutlined />}
                            loading={startLoading}
                            onClick={() => { setPartAngleName('radian1'); angleForm.submit() }}
                        >
                            设置
                        </Button></Col>
                </Row>
                <Row gutter={24}>
                    <Col span={18}>
                        <Form.Item
                            label="弧度模组2号"
                            name={['radian2', 'range']}
                            rules={[{ required: true }]}
                            tooltip="运行范围【当前弧形模组1号位置加2度】至90°"
                        >
                            <InputNumber min={-90} max={90} addonAfter="°" />
                        </Form.Item>
                    </Col>
                    <Col span={6}>
                        <Button
                            key="extra-play"
                            icon={<PlaySquareOutlined />}
                            loading={startLoading}
                            onClick={() => { setPartAngleName('radian2'); angleForm.submit() }}
                        >
                            设置
                        </Button>
                    </Col>
                </Row>
                <Row gutter={24}>
                    <Col span={18}>
                        <Form.Item
                            label="分度盘"
                            name={['dividedCircle', 'range']}
                            rules={[{ required: true }]}
                            tooltip=" 范围:±180°"
                        >
                            <InputNumber min={0} max={360} addonAfter="°" />
                        </Form.Item>
                    </Col>
                    <Col span={6}>
                        <Button
                            key="extra-play"
                            icon={<PlaySquareOutlined />}
                            loading={startLoading}
                            onClick={() => { setPartAngleName('dividedCircle'); angleForm.submit() }}
                        >
                            设置
                        </Button>
                    </Col>
                </Row>
            </Form>
        </Card>
    )
}
const SpeedConfig: React.FC = () => {

    const [setSpeedLoading, setSetSpeedLoading] = useState(false);
    const [speedForm] = Form.useForm();

    const onFinish1 = (values: any) => {
        if (setSpeedLoading) {
            setSpeed(values)
        } else {
            console.warn("no next step!");
        }
    };

    const setSpeed = (values) => {
        darkroomModuleService.setSpeed(values)
            .then(function (data) {
                notification.success({
                    message: `操作成功`,
                    description: '设置速度成功。',
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
            .finally(() => { setSetSpeedLoading(false) })
    };
    return (
        <Card title="运行速度" size='small' extra={[
            <Button
                key="extra-play"
                icon={<PlaySquareOutlined />}
                loading={setSpeedLoading}
                onClick={() => { setSetSpeedLoading(true); speedForm.submit() }}
            >
                设置全部
            </Button>
        ]}>
            <Form name="speed-form"
                form={speedForm}
                labelCol={{ span: 8 }}
                wrapperCol={{ span: 16 }}
                onFinish={onFinish1}
                initialValues={defaultValue}
            >
                <Form.Item label="直线模组1号"
                    name={['line1', 'maxSpeed']}
                    rules={[{ required: true }]}
                    tooltip=" 范围:[0,100]"
                >
                    <InputNumber min={0} max={100} addonAfter="mm/s" />
                </Form.Item>
                <Form.Item label="直线模组2号"
                    name={['line2', 'maxSpeed']}
                    rules={[{ required: true }]}
                    tooltip=" 范围:[0,100]"
                >
                    <InputNumber min={0} max={100} addonAfter="mm/s" />
                </Form.Item>
                <Form.Item label="弧形模组1号"
                    name={['radian1', 'maxSpeed']}
                    rules={[{ required: true }]}
                    tooltip=" 范围:[0,1]"
                >
                    <InputNumber min={0} max={1} addonAfter="°/s" />
                </Form.Item>
                <Form.Item label="弧形模组2号"
                    name={['radian2', 'maxSpeed']}
                    rules={[{ required: true }]}
                    tooltip=" 范围:[0,1]"
                >
                    <InputNumber min={0} max={1} addonAfter="°/s" />
                </Form.Item>
                <Form.Item label="分度盘"
                    name={['dividedCircle', 'maxSpeed']}
                    rules={[{ required: true }]}
                    tooltip=" 范围:[0,10]"
                >
                    <InputNumber min={0} max={10} addonAfter="°/s" />
                </Form.Item>
            </Form>

        </Card>
    )
}

const DarkroomModule: React.FC = () => {
    const [workStatus, setWorkStatus] = useState('unkonw');
    const [data, setData] = useState(defaultData);
    const [startLoading, setStartLoading] = useState(false);
    const [endLoading, setEndLoading] = useState(false);

    useEffect(() => {
        getStatus()
        const timeoutID = setInterval(
            getStatus
            , 2000);
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
        darkroomModuleService.start()
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
                    message: `操作失败`,
                    description: `${error}`,
                });
            })
            .finally(() => { setStartLoading(false) })
    }, [startLoading]);

    const getStatus = () => {
        darkroomModuleService.getStatus()
        darkroomModuleService.getStatusAck().then(function (data) {
            setData(data)
            setWorkStatus(data.workStatus)
        }).catch(function (error) {
            console.error(error);
            notification.error({
                message: `操作失败`,
                description: `${error}`,
            });
        })
    };

    const stop = () => {
        darkroomModuleService.stop()
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
                    message: `操作失败`,
                    description: `${error}`,
                });
            })
            .finally(() => { setEndLoading(false) })
    };

    const getExtraBtns = (workStatus) => {
        const extraBtns = [
            // <Button
            //     key="extra-download"
            //     icon={<DownloadOutlined />}
            //     loading={statusLoading}
            //     onClick={() => setStatusLoading(true)}
            // >
            //     获取状态
            // </Button>
        ]
        extraBtns.push(<Button
            key="extra-play"
            icon={<PlaySquareOutlined />}
            loading={startLoading}
            onClick={() => { setStartLoading(true); }}
        >
            开始运行
        </Button>)
        extraBtns.push(<Button
            key="extra-stop"
            icon={<StopOutlined />}
            loading={endLoading}
            onClick={() => setEndLoading(true)}
        >
            停止运行
        </Button>)
        return extraBtns
    }

    const extraBtns = getExtraBtns(workStatus);

    return (
        <>
            <PageHeader
                className="site-page-header"
                title='暗室模组'
                subTitle={(<WorkStatus workStatus={workStatus} />)}
                extra={extraBtns}
            />
            <Row gutter={24}>
                <Col offset={2} span={8}>
                    <AngleConfig />
                </Col>
                <Col span={8}>
                    <SpeedConfig />
                </Col>
            </Row>
            <Row gutter={24}>
                <Col offset={2} span={16}>
                    <hr></hr>
                </Col>
                <Col offset={2} span={16}>
                    <ModuleStatus data={data} />
                </Col>
            </Row>
        </>
    )
};

export default DarkroomModule;

