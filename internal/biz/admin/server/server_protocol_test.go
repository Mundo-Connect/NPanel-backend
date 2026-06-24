package server

import (
	"context"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	servermodel "github.com/npanel-dev/NPanel-backend/internal/model/server"
)

func TestAdminMxNormalizationDefaultsMc1AndMundoRdp(t *testing.T) {
	mc1 := &servermodel.Protocol{Type: "mx", Transport: "mc1", Security: "tls"}
	normalizeMxProtocol(mc1, "node.example.com")
	if mc1.Path != "/" {
		t.Fatalf("mc1 path = %q, want /", mc1.Path)
	}
	if mc1.Host != "node.example.com" {
		t.Fatalf("mc1 host = %q, want node.example.com", mc1.Host)
	}
	if mc1.Mc1Mode != "auto" {
		t.Fatalf("mc1 mode = %q, want auto", mc1.Mc1Mode)
	}
	if mc1.Fingerprint != "chrome" {
		t.Fatalf("mc1 fingerprint = %q, want chrome", mc1.Fingerprint)
	}

	rdp := &servermodel.Protocol{Type: "mx", Transport: "mundordp"}
	normalizeMundoProtocol(rdp)
	if rdp.MundoUsername != "MundoUser" {
		t.Fatalf("mundordp username = %q, want MundoUser", rdp.MundoUsername)
	}
}

func TestAdminCreateServerAllowsMxMc1RdpAndSql(t *testing.T) {
	repo := &adminServerRepoStub{}
	uc := NewServerUsecase(repo, &adminNodeRepoStub{}, log.DefaultLogger)

	server, err := uc.CreateServer(context.Background(), "Mundo", "US", "LA", "node.example.com", 1, []*servermodel.Protocol{
		{Type: "mx", Port: 443, Enable: true, Transport: "mc1", Security: "tls"},
		{Type: "mx", Port: 3389, Enable: true, Transport: "mundordp"},
		{Type: "mx", Port: 3306, Enable: true, Transport: "mundosql", MundoUsername: "sqluser"},
	})
	if err != nil {
		t.Fatalf("CreateServer returned error: %v", err)
	}
	if server != repo.created {
		t.Fatal("CreateServer did not return repository result")
	}
	if len(server.Protocols) != 3 {
		t.Fatalf("protocol count = %d, want 3", len(server.Protocols))
	}
	if server.Protocols[0].Transport != "mc1" ||
		server.Protocols[0].Host != "node.example.com" ||
		server.Protocols[0].Path != "/" ||
		server.Protocols[0].Mc1Mode != "auto" ||
		server.Protocols[0].Fingerprint != "chrome" {
		t.Fatalf("mc1 defaults not applied: %+v", server.Protocols[0])
	}
	if server.Protocols[1].Transport != "mundordp" || server.Protocols[1].MundoUsername != "MundoUser" {
		t.Fatalf("mundordp defaults not applied: %+v", server.Protocols[1])
	}
	if server.Protocols[2].Transport != "mundosql" || server.Protocols[2].MundoUsername != "sqluser" {
		t.Fatalf("mundosql fields not preserved: %+v", server.Protocols[2])
	}
}

type adminServerRepoStub struct {
	created *Server
}

func (r *adminServerRepoStub) CreateServer(ctx context.Context, server *Server) (*Server, error) {
	server.ID = 1
	r.created = server
	return server, nil
}

func (r *adminServerRepoStub) UpdateServer(ctx context.Context, server *Server) (*Server, error) {
	return server, nil
}

func (r *adminServerRepoStub) DeleteServer(ctx context.Context, id int) error {
	return nil
}

func (r *adminServerRepoStub) GetServerByID(ctx context.Context, id int) (*Server, error) {
	return &Server{ID: int64(id), Address: "node.example.com"}, nil
}

func (r *adminServerRepoStub) FilterServerList(ctx context.Context, page, size int32, search string) (int32, []*Server, error) {
	return 0, nil, nil
}

func (r *adminServerRepoStub) GetServerProtocols(ctx context.Context, id int) ([]*servermodel.Protocol, error) {
	return nil, nil
}

func (r *adminServerRepoStub) ResetServerSort(ctx context.Context, sortItems []*SortItem) error {
	return nil
}

func (r *adminServerRepoStub) GetServerStatus(ctx context.Context, serverID int) (*ServerResourceStatus, error) {
	return nil, nil
}

func (r *adminServerRepoStub) GetOnlineUsers(ctx context.Context, serverID int64, protocol string) (map[int64][]string, error) {
	return nil, nil
}

func (r *adminServerRepoStub) GetOnlineUsersByServer(ctx context.Context, serverID int64) (map[string]map[int64][]string, error) {
	return nil, nil
}

func (r *adminServerRepoStub) GetUserSubscribeInfo(ctx context.Context, subscribeID int) (*UserSubscribeInfo, error) {
	return nil, nil
}

type adminNodeRepoStub struct{}

func (r *adminNodeRepoStub) CreateNode(ctx context.Context, node *Node) (*Node, error) {
	return node, nil
}

func (r *adminNodeRepoStub) UpdateNode(ctx context.Context, node *Node) (*Node, error) {
	return node, nil
}

func (r *adminNodeRepoStub) DeleteNode(ctx context.Context, id int) error {
	return nil
}

func (r *adminNodeRepoStub) FilterNodeList(ctx context.Context, page, size int32, search string, nodeGroupID *int64) (int32, []*Node, error) {
	return 0, nil, nil
}

func (r *adminNodeRepoStub) ToggleNodeStatus(ctx context.Context, id int, enable *bool) (*Node, error) {
	return &Node{ID: int64(id), Enabled: enable}, nil
}

func (r *adminNodeRepoStub) QueryNodeTags(ctx context.Context) ([]string, error) {
	return nil, nil
}

func (r *adminNodeRepoStub) ResetNodeSort(ctx context.Context, sortItems []*SortItem) error {
	return nil
}

func (r *adminNodeRepoStub) ClearNodeCache(ctx context.Context, serverIDs []int) error {
	return nil
}
