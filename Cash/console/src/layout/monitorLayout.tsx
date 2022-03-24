import React, {createRef} from "react";
import {Button, Row, Col, Input, InputRef, Layout, Select} from "antd";
import Monitor from "../components/monitor";
import {get} from "../utils/ajax";
import {SyncOutlined} from '@ant-design/icons';

require("../css/console.css");
const {Option} = Select;
const {Header, Content, Footer} = Layout;
export default class MonitorLayout extends React.Component<any, any> {
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

    render(): React.ReactElement {
        return <Layout>
            <Header style={{padding: 0}}><p style={{color: "white", fontSize: "2em"}}>Cash</p></Header>
            <Content style={{margin: '24px 16px 0', height: 1125}}>
                <Monitor ref={this.monitor} addr={this.state.addr} namespace={this.state.monitorGroup}/>
            </Content>
            <Footer style={{textAlign: 'center'}}>
                {this.selector()}
            </Footer>
        </Layout>
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
            }, () => {
                if (this.monitor.current !== null) {
                    this.monitor.current.fetchGroupInfo(groups)
                }
            });
        }, () => {
            this.setState({
                groups: Array<string>()
            });
        })
    };

    selectMonitor = (value: string) => {
        let monitor = this.monitor.current;
        if (monitor !== null) {
            monitor.reset()
        }
        this.setState({
            monitorGroup: value
        });
    };

    selector = () => (
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
                <Button onClick={this.synchronize}>Synchronize<SyncOutlined/></Button>
            </Col>
        </Row>
    );
}