let Appbar = mui.react.Appbar;
let Button = mui.react.Button;
let Container = mui.react.Container;
let Row = mui.react.Row;
let Col = mui.react.Col;
let Input = mui.react.Input;

class CraigslistQueryColumnUserInput extends React.Component {
    constructor(props) {
      super(props);
      this.buttonClicked = this.buttonClicked.bind(this);
    }
    buttonClicked(){

        console.log("button clicked")
        this.props.buttonClicked()
    }
  render() {
      return (
          <div>
            <Container fluid={true}>
                <Input placeholder="Craigslist Search URL" />
                <Button onClick={this.buttonClicked} size="small" color="primary">Load</Button>
                <Button size="small" color="primary">Save</Button>
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
    state = {
        message: "initial state",
        users: [],
        error: null
    }
    constructor(props) {
      super(props);
      this.handleChange = this.handleChange.bind(this);
      this.doRequest = this.doRequest.bind(this);
    }
    handleChange(e) {
        this.props.handleChange(e.target.value);
    }
    doRequest(e) {
        console.log("dorequest)")
        const myURL="https://detroit.craigslist.org/d/architect-engineer-cad/search/egr"

        var myObj = {
            searchURL: encodeURIComponent(myURL),
        };

        fetch("http://localhost:8080/api/" , {
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            method: "POST",
            body: JSON.stringify(myObj)
        })
        .then(response => response.json())
        .then(data =>
            this.setState({
                users: data,
                craigslistQueryResponse: data.response,
                message: "ok!!",
          })
        )
        .catch(error => this.setState({ error: JSON.stringify(error), message: "something bad happened"+JSON.stringify(error.message) }))
        ;
    }

    buttonClicked(e) {
        console.log("column button clicked" + e)
        this.doRequest()
    }

    componentDidMount() {
        this.doRequest()
    }

    render() {
        return (
            <div>
                <div className="mui--text-left">
                    <CraigslistQueryColumnUserInput 
                        buttonClicked={this.buttonClicked}
                        handleChange={this.handleChange}
                        doRequest={this.doRequest}
                        />
                </div>
                    <div> 
                        The Message is: {JSON.stringify(this.state.message)}
                    </div>
                    <CraigslistQueryColumnResults 
                        results={this.state.craigslistQueryResponse}
                        hello="hello" 
                        myprop="<b>zwatef</b>"
                    />
            </div>
        );
    }
}

class App extends React.Component {
    handleQueryInput(e) {
        console.log("hello"+e)
    }
    doRequest(){
        console.log("app do request")
    }
    render() {
        return (
            <div>
                <h1>This is the App</h1>
                <hr/>
                <div className="search-table-outter wrapper">
                    <table className="search-table inner">
                        <tbody>
                    <tr>
                <td>
            <div>
                <CraigslistQueryColumn 
                    handleChange={this.handleQueryInput} 
                    doRequest={this.doRequest} 
                    />
            </div>
                </td> 
                    </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        );
    }
}

ReactDOM.render(<App />, document.getElementById('app'));

