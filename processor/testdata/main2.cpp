#include<fstream>

using namespace std;

int main() {
    ifstream fin("/tmp/a.in");
    ofstream fout("/tmp/a.out");
    int a, b;
    fin >> a >> b;
    fout << a + b << endl;
    return 0;
}