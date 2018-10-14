class CraigslistResults extends React.Component {
  render() {
    return (
    <div>{this.props.hello}</div>
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
        <div> {JSON.stringify(this.state.users)} and , {JSON.stringify(this.state.message)}</div>
        <CraigslistResults hello="hello" />
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

