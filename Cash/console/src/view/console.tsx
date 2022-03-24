import React from "react";
import {Layout, Menu} from 'antd';
import {MonitorOutlined, PlaySquareOutlined} from '@ant-design/icons';
import MonitorLayout from "../layout/monitorLayout";
import OperationPlatformLayout from "../layout/operationPlatformLayout";

const {Sider} = Layout;

export default class ConsoleView extends React.Component<any, any> {
    state = {
        currentSection: 0
    };

    changeSection = (item: any) => {
        this.setState({
            currentSection: Number(item.key)
        })
    };

    private layouts = [<MonitorLayout/>, <OperationPlatformLayout/>];

    render(): React.ReactElement {
        return <Layout>
            <Sider>
                <div className="logo"/>
                <Menu theme="dark" mode="inline" defaultSelectedKeys={['0']} onSelect={this.changeSection}>
                    <Menu.Item key="0" icon={<MonitorOutlined/>}>
                        Real Time Monitor
                    </Menu.Item>
                    <Menu.Item key="1" icon={<PlaySquareOutlined/>}>
                        Operating Platform
                    </Menu.Item>
                </Menu>
            </Sider>
            {this.layouts[this.state.currentSection]}
        </Layout>
    }
}
