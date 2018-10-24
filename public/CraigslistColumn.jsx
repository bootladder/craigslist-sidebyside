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
        this.props.doRequest(this.input.controlEl.value)
    }
    saveButtonClicked(){
        console.log("savebuttonclicked userinput" + this.props.hello)
        alert("1" + this.input.controlEl.value);
    }

    componentDidMount(){
        this.input.controlEl.value = this.props.url
        this.props.doRequest(this.props.url)
        console.log("column doRequest with: "+this.props.url)
    }
    render() {
        console.log("render CraigslistQueryUserInput: " + JSON.stringify(this.props))
        return (
<div>
    <Container>
    <Row>
        <Input className="mui--text-caption input-100percent" ref={el => { this.input = el; }} placeholder="Craigslist Search URL" />
    </Row>
    <Row>
        <Input className="input-100percent" placeholder="Search Query" />
    </Row>
    <Row>
        <select id="catAbb">
            <option value="ccc">community</option>
            <option value="eee">events</option>
            <option value="sss">for sale</option>
            <option value="ggg">gigs</option>
            <option value="hhh">housing</option>
            <option value="jjj" selected="">jobs</option>
            <option value="rrr">resumes</option>
            <option value="bbb">services</option>
        </select>
        <CitySelector />
    </Row>
    <Row>
        <Col><Input placeholder="City" /></Col>
    </Row>
    <Row>
        <Button onClick={this.loadButtonClicked} size="small" color="primary">Load</Button>
        <Button onClick={this.saveButtonClicked} size="small" color="primary">Save</Button>
    </Row>
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
                <div className="hello x-scroll y-scroll">{resultRows}</div>
            </div>
        );
    }
}

class CraigslistQueryColumn extends React.Component {

    constructor(props) {
        super(props);
        this.doRequest = this.doRequest.bind(this);
    }

    state = {
        columnState: "columnState",
        queryResponseData: "queryResponseData"
    }

    componentDidMount() {
    }

    validateCraigslistURL(url){
        console.log(JSON.stringify(url))
        if(url.length < 5){
            console.log("length too short, defaulting URL length was "+url.length)
            return "https://baltimore.craigslist.org/d/architect-engineer-cad/search/egr";
        }
        else return url;
    }
    doRequest(craigslistSearchURL){
        console.log("do request" + craigslistSearchURL + " state is " + JSON.stringify(this.state))
        var validatedURL = this.validateCraigslistURL(craigslistSearchURL)
        console.log("validatedURL: " + this.validatedURL)

        var myJsonRequestObj = {
            searchURL: encodeURIComponent(craigslistSearchURL),
        };

        fetch("http://localhost:8080/api/" , {
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            method: "POST",
            body: JSON.stringify(myJsonRequestObj)
        })
        .then(response => response.json())
        .then(data =>
            this.setState({
                users: data,
                queryResponseData: data.response,
                message: "ok!!",
                dummystate: "ok finally changed state up here!!",
          })
        )
        .catch(error => this.setState({ error: JSON.stringify(error), message: "something bad happened"+JSON.stringify(error.message) }))
        ;
    }

    render() {
        console.log("render CraigslistQueryColumn: " + JSON.stringify(this.props))
        return (
            <div>
                <div className="mui--text-left">
                    <CraigslistQueryColumnUserInput 
                        doRequest={this.doRequest}
                        url={this.props.url}
                        />
                </div>
                <CraigslistQueryColumnResults 
                    results={this.state.queryResponseData}
                    hello="hello" 
                    myprop="<b>zwatef</b>"
                />
            </div>
        );
    }
}