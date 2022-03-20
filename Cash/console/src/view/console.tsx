import React, {createRef, ReactNode} from "react";
import {Button, Layout, Menu, Row, Select, Col, Input, InputRef, Card} from 'antd';
import {UploadOutlined, UserOutlined, VideoCameraOutlined} from '@ant-design/icons';
import Monitor from "../components/monitor";
import {get} from "../utils/ajax";

const {Option} = Select;
const {Header, Content, Footer, Sider} = Layout;

export default class ConsoleView extends React.Component<any, any> {
    state = {
        groups: Array<string>(),
        addr: '',
        monitorGroup: ''
    };
    private readonly addressInputRef: React.RefObject<InputRef>;
    private readonly monitor: React.RefObject<Monitor>;

    constructor(props: any) {
        super(props);
        this.addressInputRef = createRef();
        this.monitor = createRef();
    }

    synchronize = () => {
        let addrInput = this.addressInputRef.current;
        let addr = "localhost:8000";
        if (addrInput != null && addrInput.input != null) {
            addr = addrInput.input.value;
        }

        this.setState({
            addr: addr
        });

        const url = `http://${addr}/__cash__/__groups__`;
        get(url, (groups: Array<string>) => {
            this.setState({
                groups: [...groups]
            });
        }, (error: any) => {
            console.log(error);
            this.setState({
                groups: []
            });
        })
    };

    selectMonitor = (value: string) => {
        let monitor = this.monitor.current;
        if (monitor != null) {
            monitor.reset()
        }
        this.setState({
            monitorGroup: value
        });
    };

    render(): React.ReactElement<any, string | React.JSXElementConstructor<any>> | string | number | {} | Iterable<React.ReactNode> | React.ReactPortal | boolean | null | undefined {
        return <Layout>
            <Sider>
                <div className="logo"/>
                <Menu theme="dark" mode="inline" defaultSelectedKeys={['1']}>
                    <Menu.Item key="1" icon={<UserOutlined/>}>
                        Cash 实时监控
                    </Menu.Item>
                </Menu>
            </Sider>
            <Layout>
                <Header style={{padding: 0}}><p style={{color: "white", fontSize: "2em"}}>Cash</p></Header>
                <Content style={{margin: '24px 16px 0'}}>
                    <Monitor ref={this.monitor} addr={this.state.addr} namespace={this.state.monitorGroup}/>
                </Content>
                <Footer style={{textAlign: 'center'}}>
                    <Row>
                        <Col>
                            <Input ref={this.addressInputRef} defaultValue="localhost:8000" addonBefore="http://"
                                   placeholder="address"/>
                        </Col>
                        <Col>
                            <Select defaultValue={"Choose a group"}
                                    style={{width: 150}} onSelect={this.selectMonitor}>
                                {this.state.groups.map(group => (
                                    <Option key={group}>{group}</Option>
                                ))}
                            </Select>
                        </Col>
                        <Col>
                            <Button onClick={this.synchronize}>同步</Button>
                        </Col>
                    </Row>
                </Footer>
            </Layout>
        </Layout>
    }
}
