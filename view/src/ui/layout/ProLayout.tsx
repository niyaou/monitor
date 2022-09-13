import ProLayout, { PageContainer } from '@ant-design/pro-layout';
import { Row, Col, Tabs, Form, Input, Button, PageHeader, message, Alert, Divider, notification, Space } from 'antd';
import {
  useLocation
} from "react-router-dom";
import defaultProps from './menu';
import { AppRouter } from '../router';
import React, { useState, useEffect, useContext } from 'react';
import baseApi from '../api/fc2/baseApi'
import {
  Link,
} from "react-router-dom";
// import { ipcRenderer } from 'electron'
const electronAPI = window.electronAPI
import { ConfigDataContextProvider, configDataContext, sessionFlag } from "../views/radar/fc2/configDataReducer";

export default () => {
  const [pathname, setPathname] = useState(useLocation().pathname);

  const { state, dispatch } = useContext(configDataContext);
  const _heartBeat = () => {
    baseApi.heartBeat().then((res) => {
      setTimeout(_heartBeat, 6000)
    }).catch((err) => {
      electronAPI.quit()
    })
  }


  // useEffect(() => {
  //   setTimeout(_heartBeat, 3000)
  // }, [])

  return (
    <div
      style={{
        height: '100vh',
      }}
    >
      <ProLayout
        title="RadarFc2"
        logo="https://gw.alipayobjects.com/mdn/rms_b5fcc5/afts/img/A*1NHAQYduQiQAAAAAAAAAAABkARQnAQ"
        navTheme='dark'
        headerRender={false}
        {...defaultProps}
        location={{
          pathname: pathname,
        }}
        menuItemRender={(item, dom) => {
          if (state.flag === sessionFlag.FLAG_START_SESSION) {
            //连续测试，无法切换页面
            return (<Space style={{ color: 'gray' }}>{item.name}</Space>)
          } else {
            item.disabled = state.flag === sessionFlag.FLAG_START_SESSION
            return (<Link onClick={() => {
              setPathname(item.path || useLocation().pathname);
            }} to={item.path}

            >{dom}</Link>
            )
          }
        }}
      >

        <AppRouter />
      </ProLayout>
    </div>
  );
};
