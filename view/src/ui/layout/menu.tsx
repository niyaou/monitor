import { CrownOutlined, SmileOutlined, BugOutlined, BugFilled } from '@ant-design/icons';
import React from 'react';

export default {
  route: {
    path: '/',
    routes: [
      // {
      //   path: '/radar/fc2/monitor',
      //   name: '监控',
      //   icon: <SmileOutlined />,
      // },
      {
        path: '/radar/fc2',
        name: '配置',
        icon: <CrownOutlined />,
        routes: [
          {
            path: '/radar/fc2/config2',
            name: '参数',
            icon: <CrownOutlined />,
          },
          {
            path: '/radar/fc2/turret',
            name: '转台',
            icon: <CrownOutlined />,
          },
          {
            path: '/radar/fc2/darkroom/module',
            name: '暗室模组',
            icon: <CrownOutlined />,
          },
          {
            path: '/radar/fc2/turrentTable',
            name: '室外转台',
            icon: <CrownOutlined />,
          },
        ],
      },
      // {
      //   path: '/v2x/v2xPanel',
      //   name: 'V2X测试面板',
      //   icon: <BugFilled />,
      // },
    ],
  },
  location: {
    pathname: '/',
  },
};