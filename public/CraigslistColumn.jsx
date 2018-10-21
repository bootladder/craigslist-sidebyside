let Appbar = mui.react.Appbar;
let Button = mui.react.Button;
let Container = mui.react.Container;
let Row = mui.react.Row;
let Col = mui.react.Col;
let Input = mui.react.Input;

class CraigslistQueryColumnUserInput extends React.Component {
    constructor(props) {
        super(props);
        this.loadButtonClicked = this.loadButtonClicked.bind(this);
        this.saveButtonClicked = this.saveButtonClicked.bind(this);
    }
    loadButtonClicked(){
        this.props.doRequest(url)
    }
    saveButtonClicked(){
        console.log("savebuttonclicked userinput" + this.props.hello)
        alert("1" + this.input.controlEl.value);
    }

    componentDidMount(){
        this.props.doRequest(this.input.controlEl.value)
        console.log("column doRequest with: "+this.props.url)
    }
    render() {
        console.log("my url is: " + this.props.url)
        return (
            <div>
                <Container fluid={true}>
                    <Input ref={el => { this.input = el; }} placeholder="Craigslist Search URL" />
                    <Button onClick={this.loadButtonClicked} size="small" color="primary">Load</Button>
                    <Button onClick={this.saveButtonClicked} size="small" color="primary">Save</Button>
                </Container>
            </div>
        );
    }
}

class CraigslistQueryColumnResults extends React.Component {

    render() {
        const html = $.parseHTML( this.props.results )
        
        var resultRows = []
        var i=0
        $(html).find(".result-row").each( function(q) {
            i++
            resultRows.push(React.createElement(
                'div', {key: i}, <div dangerouslySetInnerHTML= {{__html: this.innerHTML}} />
            ))
        })

        return (
            <div>
                <div className="hello">{resultRows}</div>
            </div>
        );
    }
}

class CraigslistQueryColumn extends React.Component {

    constructor(props) {
        super(props);
    }

    componentDidMount() {
//        this.props.doRequest("https://newyork.craigslist.org/d/architect-engineer-cad/search/egr")
    }

    render() {
        return (
            <div>
                <div className="mui--text-left">
                    <CraigslistQueryColumnUserInput 
                        doRequest={this.props.doRequest}
                        url={this.props.url}
                        />
                </div>
                <CraigslistQueryColumnResults 
                    results={this.props.queryResponseData}
                    hello="hello" 
                    myprop="<b>zwatef</b>"
                />
            </div>
        );
    }
}