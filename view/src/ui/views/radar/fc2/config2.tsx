import React, { useState, useEffect, useContext } from 'react';
import { Row, Col, Tabs, Form, Input, Button, PageHeader, message, Alert, Divider, notification, Statistic } from 'antd';
import lodash from 'lodash'
import getDefaultConfig from './config-default-data'
import { SettingOutlined, CaretRightOutlined, ForwardOutlined, BugOutlined, EllipsisOutlined, FolderOpenOutlined, SaveOutlined } from '@ant-design/icons';
import ProfileView from './config2/profile'
import Common from './config2/common'
import Ddma from './config2/ddma'
import PointSet from './config2/pointset'
import { configDataContext, actionTypes } from "./configDataReducer";

import baseApi from '../../../api/fc2/baseApi'
import configApi from '../../../api/fc2/configApi'
import { Cmd } from '../../../enum/fc2enums'
import { cloneDeep, find, findIndex } from 'lodash'
const { TabPane } = Tabs;

const electronAPI = window.electronAPI

declare type ConfigAction = 'config' | 'save' | 'tmpStorage';

let action: ConfigAction;

/**
 * 添加配置参数的倍数值
 * @param _configData 
 * @returns 
 */
export const parseDataByProtocal = (_configData) => {
  let configData = cloneDeep(_configData)
  let _chirps = configData.strRfeSetting.straChirpShapes
  _chirps = _chirps.map((_chirp) => {
    _chirp.u32CenterFrequency = _chirp.u32CenterFrequency * 1000
    _chirp.u32aTxChannelPower = _chirp.u32aTxChannelPower.map((power) => power * 100)
    return _chirp
  })
  configData.strRfeSetting.straChirpShapes = _chirps


  let _profiles = configData.strTef82xxFrameOpParam.straProfileOpParam

  _profiles = _profiles.map((_profile) => {
    _profile.u32aTxPhase = _profile.u32aTxPhase.map((phase) => phase * 1000)
    return _profile
  })
  configData.strTef82xxFrameOpParam.straProfileOpParam = _profiles


  let _mmd = configData.strTef82xxFrameOpParam.strPhaseRotator


  _mmd.u32aDdmaInitPhase = _mmd.u32aDdmaInitPhase.map((phase) => phase * 1000)
  _mmd.u32aDdmaPhaseUpdate = _mmd.u32aDdmaPhaseUpdate.map((phase) => phase * 1000)

  configData.strTef82xxFrameOpParam.strPhaseRotator = _mmd
  return configData
}

let cmdType = 0
const Config2: React.FC = () => {
  const [configData, setConfigData] = useState(getDefaultConfig());
  //已提交的组件path数组
  const [submitList] = useState([]);
  //组件初始化的时候，注册所有的组件path
  const [registFormArray] = useState([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [isSubmitStatus, setIsSubmitStatus] = useState<boolean>(false);
  const [resultStatus, setResultStatus] = useState(null);
  const [ids, setIds] = useState(0);
  const [activeTabKey, setActiveTabKey] = useState('common');
  const [statistic, setStatistic] = useState([]);
  const [ackMessage, setAckMessage] = useState('雷达未设置');


  const { state, dispatch } = useContext(configDataContext);

  const recvStatistic = () => {

    baseApi.getReq('ADCRecvStatistic').then((res) => {
      // console.log(res)
      setStatistic(res.data)
    }).finally(() => {
      setTimeout(recvStatistic, 1200)
    })
  }


  const Fc2ParamAck = () => {
    baseApi.getReq('Fc2ParamAck').then((res) => {
      const data = res.data
      let msg = `${data.enuAckCode === 1 ? 'ACK_UNDEFINED' : data.enuAckCode === 2 ? 'ACK_DECODE' : data.enuAckCode === 3 ? 'ACK_GENERAL' : 'ACK_WORK_STATE'}`
      let _ack_data = ''
      data.u32aAckData.map((ack_data) => { _ack_data = `${_ack_data} 0x${ack_data.toString(16)}` })
      setAckMessage(`${msg}:${_ack_data}`)
    }).finally(() => {
      setTimeout(Fc2ParamAck, 2300)
    })
  }


  useEffect(() => {
    recvStatistic()
    Fc2ParamAck()
  }, [])


  // useEffect(() => {
  //   console.log('---configData change:--', cloneDeep(configData))
  // }, [configData])

  const sendToRadar = () => {
    configData['u16Cmd'] = cmdType
    try {
      notification.info({
        message: `发出设置`,
        placement: 'topRight',
        duration: 1,
      });
      configApi.set(parseDataByProtocal(configData)).then((response) => {
        const data = response.data
        notification.success({
          message: `设置完成${JSON.stringify(data)}`,
          placement: 'topRight',
          duration: 3,
        });
        let msg = `${data.enuAckCode === 1 ? 'ACK_UNDEFINED' : data.enuAckCode === 2 ? 'ACK_DECODE' : data.enuAckCode === 3 ? 'ACK_GENERAL' : 'ACK_WORK_STATE'}`

        let _ack_data = ''
        data.u32aAckData.map((ack_data) => { _ack_data = `${_ack_data} 0x${ack_data.toString(16)}` })

        setAckMessage(`${msg}:${_ack_data}`)
        // if (data.enuParamCfgAck === 1) {
        //   notification.success({
        //     message: `设置完成`,
        //     placement: 'topRight',
        //     duration: 1,
        //   });
        //   setResultStatus({
        //     'type': "success",
        //     'msg': "设置成功",
        //   })
        // } else {
        //   console.error('错误:' + data.enuParamCfgAck);
        //   notification.error({
        //     message: `操作失败`,
        //     description: `错误内容： ${data.enuParamCfgAck}    ${response.msg}`,
        //   });
        //   setResultStatus({
        //     'type': "error",
        //     'msg': `${data}`,
        //   })
        // }
        setLoading(false);
      }).catch((error) => {
        console.error(error);
        notification.error({
          message: `操作失败`,
          description: `错误内容： ${error.message}`,
        });
        setLoading(false);
      })
    } catch (error) {
      console.error(error);
      notification.error({
        message: `操作失败`,
        description: `错误内容：  ${error.message}`,
      });
      setLoading(false);
    }
  }

  const setConfig = (currentCmdType) => {
    action = 'config'
    cmdType = currentCmdType
    setLoading(true);
    setIsSubmitStatus(true)
  };

  const saveConfig = () => {
    action = 'save'
    setLoading(true);
    setIsSubmitStatus(true)
  };

  const tmpStorageConfig = () => {
    action = 'tmpStorage'
    setLoading(true);
    setIsSubmitStatus(true)
  };

  const loadConfig = () => {
    electronAPI.open({ name: 'Fc2', extensions: ['fc2'] }).then(e => {
      if (e !== undefined) {
        console.info('---loadConfig data:--', JSON.parse(e))
        let packet = JSON.parse(e)
        let name = packet.filename

        baseApi.getReq(`ConfigFileSet?fileName=${name}_`).then((res) => {
          // console.log(res)

        })
        setConfigData(packet.data)
      }
    })
  };

  const setFailed = () => {
    console.warn('failed')
    setLoading(false);
    setIsSubmitStatus(false)
  }


  let resultStatusComp = null
  // resultStatus = <Alert message="Success Tips" type="success" showIcon /> 
  if (resultStatus?.type === "success") {
    resultStatusComp = <Alert message={resultStatus.msg} type="success" showIcon />
  } else if (resultStatus?.type === "error") {
    resultStatusComp = <Alert message={resultStatus.message} type="error" showIcon />
  } else {
    resultStatusComp = <Alert message='未设置' type="info" showIcon />
  }

  const directlySetValue = (value: any) => {
    lodash.merge(configData, value)
    setConfigData({ ...configData })
  }

  const onChange = (activeKey) => {
    setActiveTabKey(activeKey)
  }

  const setValue = (formName: string, value: any) => {
    submitList.push(formName)
    lodash.merge(configData, value)
    console.log("🚀 ~ file: config2.tsx ~ line 197 ~ setValue ~ configData", configData)

    if (lodash.difference(registFormArray, submitList).length === 0) {
      console.debug("获取数据结束", configData, registFormArray)
      submitList.splice(0, submitList.length)//清空数组
      setIsSubmitStatus(false)

      switch (action) {
        case 'config':
          sendToRadar()
          break;
        case 'save':
          electronAPI.save(JSON.stringify(configData), 'fc2')
          setLoading(false);
          break;
        case 'tmpStorage':
          dispatch({ type: actionTypes.FC_CONFIG, payload: configData })
          setLoading(false);
          notification.success({
            message: `暂存成功`,
            placement: 'topRight',
            duration: 1,
          });
          break;
        default:
          console.warn(`unkonw action:${action}`)
          break;
      }
    }
    else {
      console.debug("diff", configData, registFormArray)
    }
  }

  const setProfileValue = (index: number, value: any) => {
    console.log("setProfileValue:", index, value)
    lodash.merge(configData, value)
    setConfigData({ ...configData })
  }

  const registForm = (formName: string) => {
    registFormArray.push(formName)
    // console.debug(registFormArray)
  }

  const profiles = new Array(8)
  for (let index = 0; index < 8; index++) {
    const profileKey = `profile_${index + 1}`
    profiles[index] =
      <TabPane tab={profileKey} key={`profile-${profileKey}`}>
        <ProfileView formName={profileKey} index={index} data={configData} setProfileValue={setProfileValue} />
      </TabPane>
  }

  return (
    <PageHeader
      title='参数配置'
      subTitle={resultStatusComp}
      extra={
        [

          <Button key="extra-config-chirp-para" icon={<SettingOutlined />}
            loading={loading === true && cmdType === Cmd.CMD_CFG_CHIRP_PARA}
            disabled={loading === true && cmdType !== Cmd.CMD_CFG_CHIRP_PARA}
            onClick={() => { setConfig(Cmd.CMD_CFG_CHIRP_PARA) }}
          >
            配置Chirp
          </Button>,
          <Button key="extra-en-center-freq" icon={<SettingOutlined />}
            loading={loading === true && cmdType === Cmd.CMD_EN_CENTER_FREQ}
            disabled={loading === true && cmdType !== Cmd.CMD_EN_CENTER_FREQ}
            onClick={() => setConfig(Cmd.CMD_EN_CENTER_FREQ)}> 中心频率
          </Button>,
          <Button key="extra-start-singal-acq" icon={<CaretRightOutlined />}
            loading={loading === true && cmdType === Cmd.CMD_START_SINGAL_ACQ}
            disabled={loading === true && cmdType !== Cmd.CMD_START_SINGAL_ACQ}
            onClick={() => setConfig(Cmd.CMD_START_SINGAL_ACQ)}>单采
          </Button>,
          <Button key="extra-start-repeat-acq" icon={<ForwardOutlined />}
            loading={loading === true && cmdType === Cmd.CMD_START_REPEAT_ACQ}
            disabled={loading === true && cmdType !== Cmd.CMD_START_REPEAT_ACQ}
            onClick={() => setConfig(Cmd.CMD_START_REPEAT_ACQ)}>复采
          </Button>,
          <Button key="extra-start-continue-chirp"
            loading={loading === true && cmdType === Cmd.CMD_START_CONTINUE_CHIRP}
            disabled={loading === true && cmdType !== Cmd.CMD_START_CONTINUE_CHIRP}
            onClick={() => setConfig(Cmd.CMD_START_CONTINUE_CHIRP)}>连续发波
          </Button>,
          <Button key="extra-start-continue-chirp-second" icon={<EllipsisOutlined />}
            loading={loading === true && cmdType === Cmd.CMD_STOP_CHIRP}
            disabled={loading === true && cmdType !== Cmd.CMD_STOP_CHIRP}
            onClick={() => setConfig(Cmd.CMD_STOP_CHIRP)}>结束发波
          </Button>,
          <Divider type="vertical" key="extra-divider1" />,

          <Button
            key="extra-local-load"
            icon={<FolderOpenOutlined />}
            onClick={() => loadConfig()}
            disabled={activeTabKey === 'common' || activeTabKey === 'ddma' || activeTabKey === 'point' ? false : true}
          >
            加载配置
          </Button>,
          <Button
            key="extra-local-save"
            icon={<SaveOutlined />}
            onClick={() => saveConfig()}
            disabled={activeTabKey === 'common' || activeTabKey === 'ddma' || activeTabKey === 'point' ? false : true}
          >
            保存配置
          </Button>,
          <Button
            key="extra-tmp-storage"
            icon={<SaveOutlined />}
            onClick={() => tmpStorageConfig()}
            disabled={activeTabKey === 'common' || activeTabKey === 'ddma' || activeTabKey === 'point' ? false : true}
          >
            暂存配置
          </Button>,
        ]}
    >
      <Row justify="space-between" align="middle" gutter={12}>
        <Col span={12}>
          <Alert message={ackMessage} type="info" showIcon />
        </Col>
        <Col span={10} >
          <Row justify="start" gutter={12}>
            <Col span={8} >
              <Statistic title="saved Len" value={statistic[2]} />
            </Col>
            <Col span={8}>
              <Statistic title="send Len" value={statistic[1]} />
            </Col>
            <Col span={8}>
              <Statistic title="frames" value={statistic[0]} />
            </Col>
          </Row>
        </Col>

      </Row>

      <Tabs defaultActiveKey="common" destroyInactiveTabPane={true} onChange={onChange}>
        <TabPane tab="通用配置" key="common" forceRender={true}>
          <Common formName="通用配置" data={configData} isSubmitStatus={isSubmitStatus} setValue={setValue} directlySetValue={directlySetValue} setFailed={setFailed} registForm={registForm} />
        </TabPane>
        <TabPane tab="DDMA控制" key="ddma" forceRender={true}>
          <Ddma formName="DDMA控制" data={configData} isSubmitStatus={isSubmitStatus} setValue={setValue} directlySetValue={directlySetValue} setFailed={setFailed} registForm={registForm} />
        </TabPane>
        {profiles}
        <TabPane tab="点频配置" key="point" forceRender={true}>
          <PointSet formName="点频配置" index={0} data={configData} setProfileValue={setProfileValue} isSubmitStatus={isSubmitStatus} setValue={setValue} setFailed={setFailed} registForm={registForm} />
        </TabPane>
      </Tabs>
    </PageHeader >
  )
};

export default Config2;