class CraigslistQueryColumnResults extends React.Component {

  render() {
    const html = $.parseHTML( this.props.results ) //console.log(html.html)
    
    //const b = $(html).find(".result-row").length; console.log(b)
    //const z = $(html).find(".result-row").html()
      
    var resultRows = []
    var i=0
    $(html).find(".result-row").each( function(q) {
        //console.log(this.innerHTML)
        i++
        resultRows.push(React.createElement(
            'div', {key: i}, <div dangerouslySetInnerHTML= {{__html: this.innerHTML}} />
        ))
    })

    var stations = [];
    var i
    for(i=0; i<4; i++){
      stations.push(
          <div key={i} className="station">
              Call: {i}, Freq: {i} <br/>
          </div>
      )
    }

    return (
      <div>
            <div className="stations">{stations}</div>
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
  }
  handleChange(e) {
    this.props.onInputEvent(e.target.value);
  }

  componentDidMount() {

    const myURL="https://detroit.craigslist.org/d/architect-engineer-cad/search/egr"

    var myObj = {
      searchURL: encodeURIComponent(myURL),
      lastName:"Doe", age:50, eyeColor:"blue"
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
    .catch(error => this.setState({ error: JSON.stringify(error), message: "something bad happened"+JSON.stringify(error.message) }));
  }

  render() {
    return (
      <div>
        <form action="/action_page.php">
          URL: <input type="text" name="url"/><br/>
        </form> 
        <button />
        <input value={37} onChange={this.handleChange} />
        <div> 
          {JSON.stringify(this.state.message)}
        </div>
        <CraigslistQueryColumnResults hello="hello" results={this.state.craigslistQueryResponse}
              myprop="<b>zwatef</b>"/>
      </div>
    );
  }
}

class App extends React.Component {
  handleQueryInput(e) {
    console.log("hello"+e)
  }
  render() {
    return (
      <div>
        <h1>Hello JSX and the World!</h1>
        <CraigslistQueryColumn onInputEvent={this.handleQueryInput}/>
      </div>
    );
  }
}

ReactDOM.render(<App />, document.getElementById('app'));

